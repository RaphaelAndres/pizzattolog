package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ─── Enums ───────────────────────────────────────────────────────────────────

type Role string
type StatusLicenca string
type TipoLicenca string

const (
	RoleAdmin       Role = "admin"
	RoleGestor      Role = "gestor"
	RoleVisualizador Role = "visualizador"
)

const (
	StatusAtiva             StatusLicenca = "ativa"
	StatusProximaVencimento StatusLicenca = "proxima_vencimento"
	StatusVencida           StatusLicenca = "vencida"
)

const (
	TipoAmbiental  TipoLicenca = "ambiental"
	TipoPoliciaC   TipoLicenca = "policia_civil"
	TipoSanitaria  TipoLicenca = "sanitaria"
	TipoBombeiros  TipoLicenca = "bombeiros"
	TipoPrefeitura TipoLicenca = "prefeitura"
	TipoOutro      TipoLicenca = "outro"
)

// ─── Models ──────────────────────────────────────────────────────────────────

type Usuario struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Nome         string         `gorm:"size:150;not null" json:"nome"`
	Email        string         `gorm:"size:200;uniqueIndex;not null" json:"email"`
	SenhaHash    string         `gorm:"size:255;not null" json:"-"`
	Role         Role           `gorm:"size:20;default:'gestor'" json:"role"`
	Ativo        bool           `gorm:"default:true" json:"ativo"`
	CreatedAt    time.Time      `json:"criado_em"`
	UpdatedAt    time.Time      `json:"atualizado_em"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *Usuario) SetSenha(senha string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.SenhaHash = string(hash)
	return nil
}

func (u *Usuario) VerificarSenha(senha string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.SenhaHash), []byte(senha)) == nil
}

type Licenca struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Nome            string         `gorm:"size:200;not null" json:"nome"`
	Tipo            TipoLicenca    `gorm:"size:30;not null" json:"tipo"`
	OrgaoEmissor    string         `gorm:"size:150" json:"orgao_emissor"`
	Numero          string         `gorm:"size:100" json:"numero"`
	Descricao       string         `gorm:"type:text" json:"descricao"`
	DataEmissao     *time.Time     `json:"data_emissao"`
	DataValidade    time.Time      `gorm:"not null;index" json:"data_validade"`
	Status          StatusLicenca  `gorm:"size:30;default:'ativa'" json:"status"`
	ArquivoKey      string         `gorm:"size:500" json:"arquivo_key"`
	ArquivoNome     string         `gorm:"size:255" json:"arquivo_nome"`
	ArquivoTamanho  int64          `json:"arquivo_tamanho"`
	CriadoPorID     uint           `json:"criado_por_id"`
	CriadoPor       *Usuario       `gorm:"foreignKey:CriadoPorID" json:"criado_por,omitempty"`
	CreatedAt       time.Time      `json:"criado_em"`
	UpdatedAt       time.Time      `json:"atualizado_em"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// DiasParaVencer retorna quantos dias faltam para vencer (negativo = vencida)
func (l *Licenca) DiasParaVencer() int {
	return int(time.Until(l.DataValidade).Hours() / 24)
}

// AtualizarStatus atualiza o status baseado na data de validade
func (l *Licenca) AtualizarStatus() {
	dias := l.DiasParaVencer()
	switch {
	case dias < 0:
		l.Status = StatusVencida
	case dias <= 30:
		l.Status = StatusProximaVencimento
	default:
		l.Status = StatusAtiva
	}
}

type AlertaEnviado struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	LicencaID    uint      `gorm:"index;not null" json:"licenca_id"`
	Licenca      *Licenca  `gorm:"foreignKey:LicencaID" json:"licenca,omitempty"`
	TipoAlerta   string    `gorm:"size:10;not null" json:"tipo_alerta"` // "30d", "15d", "7d"
	EnviadoEm    time.Time `json:"enviado_em"`
	Destinatarios string   `gorm:"type:json" json:"destinatarios"`
}

// ─── Conexão ─────────────────────────────────────────────────────────────────

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	logLevel := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Usuario{}, &Licenca{}, &AlertaEnviado{})
}

func SeedAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&Usuario{}).Where("role = ?", RoleAdmin).Count(&count)
	if count > 0 {
		return nil
	}

	admin := &Usuario{
		Nome:  "Administrador",
		Email: "admin@pizzattolog.com.br",
		Role:  RoleAdmin,
		Ativo: true,
	}
	if err := admin.SetSenha("Admin@123"); err != nil {
		return err
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	log.Println("✅ Usuário admin criado: admin@pizzattolog.com.br / Admin@123")
	return nil
}
