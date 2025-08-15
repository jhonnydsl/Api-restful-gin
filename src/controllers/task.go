package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/api-restful-gin/src/dtos"
	dtosPage "github.com/jhonnydsl/api-restful-gin/src/dtos/pagination"
	"github.com/jhonnydsl/api-restful-gin/src/entities"
	"github.com/jhonnydsl/api-restful-gin/src/repositorys"
	"github.com/jhonnydsl/api-restful-gin/src/services"
	"github.com/jhonnydsl/api-restful-gin/src/utils"
	"github.com/jhonnydsl/api-restful-gin/src/utils/converts"
	"github.com/jhonnydsl/api-restful-gin/src/utils/enum"
	"github.com/jhonnydsl/api-restful-gin/src/utils/formats"
)

type TaskController struct {
	service *services.TaskService
}

func NewTaskController(server *gin.Engine, repo *repositorys.TaskRepository) {
	service := services.NewTaskService(repo)
	controller := &TaskController{
		service: service,
	}

	routes := server.Group("/tasks")
	{
		routes.POST("", controller.CreateTask)
		routes.GET("/:id", controller.GetTaskByID)
		routes.GET("", controller.GetTasks)
		routes.PUT("/:id", controller.UpdateTask)
		routes.DELETE("/:id", controller.DeleteTask)
	}
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks [post]
// @Summary Criar uma nova tarefa
// @Description Registra uma nova tarefa na API
// @Accept json
// @Produce json
// @Param task body dtos.Task true "Dados do usuário"
// @Success 201 {object} dtos.Message "Tarefa criada"
// @Failure 400 {object} dtos.APIError "Erro de validação"
// @Failure 409 {object} dtos.APIError "Erro de conflito, dados ja existem"
func (c *TaskController) CreateTask(ginContext *gin.Context) {
	var taskDTO dtos.Task

	err := utils.ValidateRequestBody(ginContext, &taskDTO)
	if err != nil {
		ginContext.Error(err)
		return
	}

	// Recupera o ID do usuário autenticado via JWT.
	userID, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	// Verifica se já existe uma tarefa com o mesmo título para o mesmo usuário.
	err = c.service.ExistsTaskByTitle(dtos.ExistsFilter{Field: "title", Value: taskDTO.Title, ForeignKey: "userID", ForeignKeyValue: userID})
	if err != nil {
		ginContext.Error(err)
		return
	}
	
	err = c.service.CreateTask(userID, taskDTO.Title, taskDTO.Description)
	if err != nil {
		ginContext.Error(err)
		return
	}
	ginContext.JSON(http.StatusOK, dtos.Message{
		Message: "Tarefa criada com sucesso.",
	})
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks/{id} [get]
// @Summary Buscar uma nova tarefa por ID
// @Description Busca uma nova tarefa na API
// @Accept json
// @Produce json
// @Param id path string true "ID da Task" example("60c72b2f9b1d8b57b8ed2123")
// @Success 200 {object} entities.Task "Tarefa por ID"
// @Failure 400 {object} dtos.APIError "Erro de validação"
func (c *TaskController) GetTaskByID(ginContext *gin.Context) {
	var task entities.Task
	idHex := ginContext.Param("id")

	id, err := converts.StringToObject(idHex)
	if err != nil {
		ginContext.Error(err)
		return
	}

	userID, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	param := dtos.GetAnyFilter{
		Field: "_id",
		Value: id,
		ForeignKey: "userID",
		ForeignKeyValue: userID,
		Result: &task,
	}

	err = c.service.GetTaskByID(param)
	if err != nil {
		ginContext.Error(err)
		return
	}

	createdAt := utils.TimeBrazil(task.CreateAt)
	updatedAt := utils.TimeBrazil(task.UpdateAt)

	task.CreateAt, err = formats.Time(createdAt, enum.FormatTime.DataHour())
	if err != nil {
		ginContext.Error(err)
		return
	}

	task.UpdateAt, err = formats.Time(updatedAt, enum.FormatTime.DataHour())
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, task)
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks [get]
// @Summary Pegar a paginação
// @Description Pegar a paginação de tarefas na API
// @Accept json
// @Produce json
// @Param page query int false "Número da pagina"
// @Param limitPage query int false "Número de registros por página"
// @Param searchField query string false "Buscar por qual propriedade de task"
// @Param searchValue query string false "Buscar pelo valor da propriedade da task"
// @Param sortField query string false "Ordenar por qual propriedade de task"
// @Param sortOrder query string false "Ordenação" Enums(ascending, descending)
// @Success 200 {array} dtosPage.PaginationResult[entities.Task] "Lista de tasks"
// @Failure 400 {object} dtos.APIError "Erro de validação"
func (c *TaskController) GetTasks(ginContext *gin.Context) {
	page := ginContext.DefaultQuery("page", "1")
	limitPage := ginContext.DefaultQuery("limitPage", "5")
	searchField := ginContext.DefaultQuery("searchField", "")
	searchValue := ginContext.DefaultQuery("searchValue", "")
	sortField := ginContext.DefaultQuery("sortField", "_id")
	sortOrderStr := ginContext.DefaultQuery("sortOrder", enum.SortOrder.AscendingStr())

	sortOrder, err := enum.SortOrder.ConvertSortOrderEnumToInt(sortOrderStr)
	if err != nil {
		ginContext.Error(err)
		return
	}

	skipInt, err := converts.StringToInt(page)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, dtos.APIError{Message: "Parâmetro skip inválido."})
		return
	}

	limitInt, err := converts.StringToInt(limitPage)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, dtos.APIError{Message: "Parâmetro perPage inválido."})
		return
	}

	userID, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}
	var tasks []entities.Task

	paginationResult := c.service.GetPagination(dtosPage.PaginationParams{
		Field: "userID",
		Value: userID,
		Result: &tasks,
		Skip: skipInt,
		Limit: limitInt,
		SearchField: searchField,
		SearchValue: searchValue,
		SortField: sortField,
		SortOrder: sortOrder,
	})

	for item := range tasks {
		tasks[item].CreateAt = utils.TimeBrazil(tasks[item].CreateAt)
		tasks[item].UpdateAt = utils.TimeBrazil(tasks[item].UpdateAt)
	}

	if paginationResult.PaginationResultContext.Err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, paginationResult)
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks/{id} [put]
// @Summary Editar uma tarefa por ID
// @Description Edita uma tarefa na API
// @Accept json
// @Produce json
// @Param id path string true "ID da Task" example("60c72b2f9b1d8b57b8ed2123")
// @Param task body dtos.Task true "Dados do usuário"
// @Success 200 {object} dtos.Message "Tarefa por ID"
// @Failure 400 {object} dtos.APIError "Erro de validação"
func (c *TaskController) UpdateTask(ginContext *gin.Context) {
	var updateData dtos.Task

	err := utils.ValidateRequestBody(ginContext, &updateData)
	if err != nil {
		ginContext.Error(err)
		return
	}

	taskIDStr := ginContext.Param("id")

	taskID, err := converts.StringToObject(taskIDStr)
	if err != nil {
		ginContext.Error(err)
		return
	}

	userID, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.UpdateTask(taskID, userID, updateData)
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, dtos.Message{
		Message: "Tarefa atualizada com sucesso.",
	})
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks/{id} [delete]
// @Summary Deletar uma tarefa por ID
// @Description Deleta uma tarefa na API
// @Accept json
// @Produce json
// @Param id path string true "ID da Task" example("60c72b2f9b1d8b57b8ed2123")
// @Success 200 {object} dtos.Message "Tarefa por ID"
// @Failure 400 {object} dtos.APIError "Erro de validação"
func (c *TaskController) DeleteTask(ginContext *gin.Context) {
	id := ginContext.Param("id")

	idObj, err := converts.StringToObject(id)
	if err != nil {
		ginContext.Error(err)
		return
	}

	userID, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	params := dtos.DeleteFilter{
		ID: idObj,
		ForeignKey: "userID",
		ForeignKeyValue: userID,
	}

	err = c.service.DeleteTask(params)
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, dtos.Message{
		Message: "Tarefa deletada com sucesso.",
	})
}