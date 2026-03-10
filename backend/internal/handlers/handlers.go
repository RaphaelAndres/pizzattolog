package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pizzattolog/licencas/internal/models"
	"github.com/pizzattolog/licencas/internal/repository"
	"github.com/pizzattolog/licencas/internal/services"
)

// ─── Auth Handler ─────────────────────────────────────────────────────────────

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(s *services.UserService) *AuthHandler {
	return &AuthHandler{userService: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Senha string `json:"senha" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	resp, err := h.userService.Login(req.Email, req.Senha)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"mensagem": "refresh não implementado no MVP"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"mensagem": "logout realizado"})
}

func (h *AuthHandler) ListUsuarios(c *gin.Context) {
	usuarios, err := h.userService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao listar usuários"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": usuarios})
}

func (h *AuthHandler) CreateUsuario(c *gin.Context) {
	var req struct {
		Nome  string      `json:"nome" binding:"required"`
		Email string      `json:"email" binding:"required,email"`
		Senha string      `json:"senha" binding:"required,min=8"`
		Role  models.Role `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	user, err := h.userService.CreateUsuario(req.Nome, req.Email, req.Senha, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"dados": user})
}

func (h *AuthHandler) UpdateUsuario(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Nome  string `json:"nome"`
		Ativo bool   `json:"ativo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	user, err := h.userService.Update(uint(id), req.Nome, req.Ativo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": user})
}

func (h *AuthHandler) DeleteUsuario(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao remover usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "Usuário removido"})
}

// ─── Licença Handler ─────────────────────────────────────────────────────────

type LicencaHandler struct {
	service *services.LicencaService
}

func NewLicencaHandler(s *services.LicencaService) *LicencaHandler {
	return &LicencaHandler{service: s}
}

func (h *LicencaHandler) List(c *gin.Context) {
	filtros := repository.LicencaFiltros{
		Tipo:   c.Query("tipo"),
		Status: c.Query("status"),
		Busca:  c.Query("busca"),
		Ordem:  c.Query("ordem"),
	}

	licencas, err := h.service.List(filtros)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao listar licenças"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": licencas, "total": len(licencas)})
}

func (h *LicencaHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	licenca, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Licença não encontrada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": licenca})
}

func (h *LicencaHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	nome := c.PostForm("nome")
	tipo := c.PostForm("tipo")
	orgao := c.PostForm("orgao_emissor")
	numero := c.PostForm("numero")
	descricao := c.PostForm("descricao")
	dataValStr := c.PostForm("data_validade")

	if nome == "" || tipo == "" || dataValStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Nome, tipo e data de validade são obrigatórios"})
		return
	}

	dataVal, err := time.Parse("2006-01-02", dataValStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato de data inválido (use YYYY-MM-DD)"})
		return
	}

	licenca := &models.Licenca{
		Nome:         nome,
		Tipo:         models.TipoLicenca(tipo),
		OrgaoEmissor: orgao,
		Numero:       numero,
		Descricao:    descricao,
		DataValidade: dataVal,
	}

	// Arquivo opcional
	var file interface{ Read([]byte) (int, error) } = nil
	_ = file
	fh, err := c.FormFile("arquivo")
	if err == nil {
		f, _ := fh.Open()
		defer f.Close()
		result, err := h.service.Create(licenca, f, fh, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"dados": result})
		return
	}

	result, err := h.service.Create(licenca, nil, nil, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"dados": result})
}

func (h *LicencaHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	dataValStr := c.PostForm("data_validade")
	dataVal, err := time.Parse("2006-01-02", dataValStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato de data inválido"})
		return
	}

	dados := &models.Licenca{
		Nome:         c.PostForm("nome"),
		Tipo:         models.TipoLicenca(c.PostForm("tipo")),
		OrgaoEmissor: c.PostForm("orgao_emissor"),
		Numero:       c.PostForm("numero"),
		Descricao:    c.PostForm("descricao"),
		DataValidade: dataVal,
	}

	fh, err := c.FormFile("arquivo")
	if err == nil {
		f, _ := fh.Open()
		defer f.Close()
		result, err := h.service.Update(uint(id), dados, f, fh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"dados": result})
		return
	}

	result, err := h.service.Update(uint(id), dados, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": result})
}

func (h *LicencaHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao remover licença"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "Licença removida com sucesso"})
}

func (h *LicencaHandler) DownloadArquivo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	url, err := h.service.GetArquivoURL(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// ─── Dashboard Handler ────────────────────────────────────────────────────────

type DashboardHandler struct {
	licencaService *services.LicencaService
	alertaService  *services.AlertaService
}

func NewDashboardHandler(s *services.LicencaService) *DashboardHandler {
	return &DashboardHandler{licencaService: s}
}

func (h *DashboardHandler) GetResumo(c *gin.Context) {
	resumo, err := h.licencaService.GetResumo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar resumo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": resumo})
}

func (h *DashboardHandler) GetAlertas(c *gin.Context) {
	filtros := repository.LicencaFiltros{Status: "proxima_vencimento"}
	licencas, err := h.licencaService.List(filtros)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar alertas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": licencas, "total": len(licencas)})
}
