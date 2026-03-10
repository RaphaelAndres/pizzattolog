package repository

import (
	"time"

	"github.com/pizzattolog/licencas/internal/models"
	"gorm.io/gorm"
)

// ─── User Repository ─────────────────────────────────────────────────────────

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.Usuario, error) {
	var user models.Usuario
	err := r.db.Where("email = ? AND ativo = true", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uint) (*models.Usuario, error) {
	var user models.Usuario
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Create(user *models.Usuario) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.Usuario) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}

func (r *UserRepository) List() ([]models.Usuario, error) {
	var users []models.Usuario
	err := r.db.Find(&users).Error
	return users, err
}

// ─── Licença Repository ──────────────────────────────────────────────────────

type LicencaRepository struct {
	db *gorm.DB
}

func NewLicencaRepository(db *gorm.DB) *LicencaRepository {
	return &LicencaRepository{db: db}
}

type LicencaFiltros struct {
	Tipo    string
	Status  string
	Busca   string
	Ordem   string
}

func (r *LicencaRepository) List(filtros LicencaFiltros) ([]models.Licenca, error) {
	var licencas []models.Licenca
	q := r.db.Preload("CriadoPor")

	if filtros.Tipo != "" {
		q = q.Where("tipo = ?", filtros.Tipo)
	}
	if filtros.Status != "" {
		q = q.Where("status = ?", filtros.Status)
	}
	if filtros.Busca != "" {
		like := "%" + filtros.Busca + "%"
		q = q.Where("nome LIKE ? OR numero LIKE ? OR orgao_emissor LIKE ?", like, like, like)
	}

	ordem := "data_validade ASC"
	if filtros.Ordem != "" {
		ordem = filtros.Ordem
	}
	q = q.Order(ordem)

	err := q.Find(&licencas).Error
	return licencas, err
}

func (r *LicencaRepository) FindByID(id uint) (*models.Licenca, error) {
	var licenca models.Licenca
	err := r.db.Preload("CriadoPor").First(&licenca, id).Error
	return &licenca, err
}

func (r *LicencaRepository) Create(licenca *models.Licenca) error {
	return r.db.Create(licenca).Error
}

func (r *LicencaRepository) Update(licenca *models.Licenca) error {
	return r.db.Save(licenca).Error
}

func (r *LicencaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Licenca{}, id).Error
}

// FindVencendoEm busca licenças que vencem em exatamente N dias (para alertas)
func (r *LicencaRepository) FindVencendoEm(dias int) ([]models.Licenca, error) {
	alvo := time.Now().AddDate(0, 0, dias)
	inicio := time.Date(alvo.Year(), alvo.Month(), alvo.Day(), 0, 0, 0, 0, alvo.Location())
	fim := inicio.Add(24 * time.Hour)

	var licencas []models.Licenca
	err := r.db.Preload("CriadoPor").
		Where("data_validade >= ? AND data_validade < ? AND status != ?", inicio, fim, models.StatusVencida).
		Find(&licencas).Error
	return licencas, err
}

// FindProximasVencer retorna licenças que vencem nos próximos N dias
func (r *LicencaRepository) FindProximasVencer(dias int) ([]models.Licenca, error) {
	limite := time.Now().AddDate(0, 0, dias)
	var licencas []models.Licenca
	err := r.db.Where("data_validade <= ? AND data_validade >= ?", limite, time.Now()).
		Order("data_validade ASC").
		Find(&licencas).Error
	return licencas, err
}

// ContarPorStatus retorna contagem por status
func (r *LicencaRepository) ContarPorStatus() (map[string]int64, error) {
	result := map[string]int64{}
	type row struct {
		Status string
		Total  int64
	}
	var rows []row
	err := r.db.Model(&models.Licenca{}).
		Select("status, count(*) as total").
		Group("status").
		Scan(&rows).Error
	for _, rw := range rows {
		result[rw.Status] = rw.Total
	}
	return result, err
}

// AlertaJaEnviado verifica se o alerta já foi disparado
func (r *LicencaRepository) AlertaJaEnviado(licencaID uint, tipo string) bool {
	var count int64
	hoje := time.Now()
	inicio := time.Date(hoje.Year(), hoje.Month(), hoje.Day(), 0, 0, 0, 0, hoje.Location())

	r.db.Model(&models.AlertaEnviado{}).
		Where("licenca_id = ? AND tipo_alerta = ? AND enviado_em >= ?", licencaID, tipo, inicio).
		Count(&count)
	return count > 0
}

// RegistrarAlerta salva registro de alerta enviado
func (r *LicencaRepository) RegistrarAlerta(alerta *models.AlertaEnviado) error {
	return r.db.Create(alerta).Error
}
