package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pizzattolog/licencas/internal/auth"
	"github.com/pizzattolog/licencas/internal/models"
	"github.com/pizzattolog/licencas/internal/repository"
	"github.com/robfig/cron/v3"
)

// ─── MinIO Service ───────────────────────────────────────────────────────────

func NewMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente MinIO: %w", err)
	}

	// Garante que o bucket existe
	bucket := os.Getenv("MINIO_BUCKET")
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("erro ao criar bucket: %w", err)
		}
	}

	return client, nil
}

// ─── User Service ─────────────────────────────────────────────────────────────

type UserService struct {
	repo       *repository.UserRepository
	jwtService *auth.JWTService
}

func NewUserService(repo *repository.UserRepository, jwt *auth.JWTService) *UserService {
	return &UserService{repo: repo, jwtService: jwt}
}

type LoginResponse struct {
	Token   string          `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

func (s *UserService) Login(email, senha string) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || !user.VerificarSenha(senha) {
		return nil, errors.New("e-mail ou senha inválidos")
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return &LoginResponse{Token: token, Usuario: user}, nil
}

func (s *UserService) CreateUsuario(nome, email, senha string, role models.Role) (*models.Usuario, error) {
	user := &models.Usuario{Nome: nome, Email: email, Role: role, Ativo: true}
	if err := user.SetSenha(senha); err != nil {
		return nil, err
	}
	return user, s.repo.Create(user)
}

func (s *UserService) List() ([]models.Usuario, error) {
	return s.repo.List()
}

func (s *UserService) Update(id uint, nome string, ativo bool) (*models.Usuario, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	user.Nome = nome
	user.Ativo = ativo
	return user, s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ─── Licença Service ─────────────────────────────────────────────────────────

type LicencaService struct {
	repo        *repository.LicencaRepository
	minioClient *minio.Client
}

func NewLicencaService(repo *repository.LicencaRepository, minioClient *minio.Client) *LicencaService {
	return &LicencaService{repo: repo, minioClient: minioClient}
}

func (s *LicencaService) List(filtros repository.LicencaFiltros) ([]models.Licenca, error) {
	return s.repo.List(filtros)
}

func (s *LicencaService) GetByID(id uint) (*models.Licenca, error) {
	return s.repo.FindByID(id)
}

func (s *LicencaService) Create(
	licenca *models.Licenca,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	userID uint,
) (*models.Licenca, error) {
	licenca.CriadoPorID = userID
	licenca.AtualizarStatus()

	if file != nil && fileHeader != nil {
		key, err := s.uploadArquivo(file, fileHeader)
		if err != nil {
			return nil, fmt.Errorf("erro no upload: %w", err)
		}
		licenca.ArquivoKey = key
		licenca.ArquivoNome = fileHeader.Filename
		licenca.ArquivoTamanho = fileHeader.Size
	}

	if err := s.repo.Create(licenca); err != nil {
		return nil, err
	}
	return licenca, nil
}

func (s *LicencaService) Update(id uint, dados *models.Licenca, file multipart.File, fileHeader *multipart.FileHeader) (*models.Licenca, error) {
	licenca, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("licença não encontrada")
	}

	licenca.Nome = dados.Nome
	licenca.Tipo = dados.Tipo
	licenca.OrgaoEmissor = dados.OrgaoEmissor
	licenca.Numero = dados.Numero
	licenca.Descricao = dados.Descricao
	licenca.DataEmissao = dados.DataEmissao
	licenca.DataValidade = dados.DataValidade
	licenca.AtualizarStatus()

	if file != nil && fileHeader != nil {
		key, err := s.uploadArquivo(file, fileHeader)
		if err != nil {
			return nil, fmt.Errorf("erro no upload: %w", err)
		}
		licenca.ArquivoKey = key
		licenca.ArquivoNome = fileHeader.Filename
		licenca.ArquivoTamanho = fileHeader.Size
	}

	return licenca, s.repo.Update(licenca)
}

func (s *LicencaService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *LicencaService) GetArquivoURL(id uint) (string, error) {
	licenca, err := s.repo.FindByID(id)
	if err != nil {
		return "", errors.New("licença não encontrada")
	}
	if licenca.ArquivoKey == "" {
		return "", errors.New("licença sem arquivo")
	}

	bucket := os.Getenv("MINIO_BUCKET")
	url, err := s.minioClient.PresignedGetObject(context.Background(), bucket, licenca.ArquivoKey, 1*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar URL: %w", err)
	}
	return url.String(), nil
}

func (s *LicencaService) uploadArquivo(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("licencas/%s%s", uuid.New().String(), ext)
	bucket := os.Getenv("MINIO_BUCKET")

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.minioClient.PutObject(
		context.Background(),
		bucket,
		key,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return key, err
}

func (s *LicencaService) GetResumo() (map[string]interface{}, error) {
	contagem, err := s.repo.ContarPorStatus()
	if err != nil {
		return nil, err
	}

	proximas, err := s.repo.FindProximasVencer(30)
	if err != nil {
		return nil, err
	}

	// Limita a 5 para o dashboard
	if len(proximas) > 5 {
		proximas = proximas[:5]
	}

	return map[string]interface{}{
		"contagem_por_status": contagem,
		"proximas_vencer":     proximas,
	}, nil
}

// ─── Alerta Service ──────────────────────────────────────────────────────────

type AlertaService struct {
	repo *repository.LicencaRepository
}

func NewAlertaService(repo *repository.LicencaRepository) *AlertaService {
	return &AlertaService{repo: repo}
}

func (s *AlertaService) StartCron() {
	c := cron.New()

	// Roda todo dia às 08:00
	c.AddFunc("0 8 * * *", func() {
		log.Println("🔔 Verificando licenças próximas ao vencimento...")
		s.verificarAlertas()
	})

	c.Start()
	log.Println("✅ Cron de alertas iniciado (08:00 diário)")
}

func (s *AlertaService) verificarAlertas() {
	for _, dias := range []int{30, 15, 7} {
		tipo := fmt.Sprintf("%dd", dias)
		licencas, err := s.repo.FindVencendoEm(dias)
		if err != nil {
			log.Printf("Erro ao buscar licenças para alerta %s: %v", tipo, err)
			continue
		}

		for _, l := range licencas {
			if s.repo.AlertaJaEnviado(l.ID, tipo) {
				continue
			}
			s.enviarAlerta(l, dias)
		}
	}
}

func (s *AlertaService) enviarAlerta(licenca models.Licenca, dias int) {
	log.Printf("📧 Alerta: Licença '%s' vence em %d dias (%s)", licenca.Nome, dias, licenca.DataValidade.Format("02/01/2006"))

	alerta := &models.AlertaEnviado{
		LicencaID:    licenca.ID,
		TipoAlerta:   fmt.Sprintf("%dd", dias),
		EnviadoEm:    time.Now(),
		Destinatarios: `[]`,
	}
	if err := s.repo.RegistrarAlerta(alerta); err != nil {
		log.Printf("Erro ao registrar alerta: %v", err)
	}
}

func (s *AlertaService) GetAlertas() ([]models.Licenca, error) {
	return s.repo.FindProximasVencer(30)
}
