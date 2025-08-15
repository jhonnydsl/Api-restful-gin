package services

import (
	"time"

	"github.com/jhonnydsl/api-restful-gin/src/dtos"
	dtosPage "github.com/jhonnydsl/api-restful-gin/src/dtos/pagination"
	"github.com/jhonnydsl/api-restful-gin/src/entities"
	"github.com/jhonnydsl/api-restful-gin/src/repositorys"
	"github.com/jhonnydsl/api-restful-gin/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo *repositorys.TaskRepository
}

func NewTaskService(repo *repositorys.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (service *TaskService) CreateTask(userID primitive.ObjectID, title, description string) error {
	task := entities.Task{
		Title: title,
		Description: description,
		UserID: userID,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	contextServer := utils.CreateContextServerWithTimeout()

	err := service.repo.Create(contextServer, task)
	if err != nil {
		return err
	}

	return nil
}

func (service *TaskService) ExistsTaskByTitle(param dtos.ExistsFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()

	err := service.repo.ExistsByAny(contextServer, param)
	if err != nil {
		return err
	}
	return nil
}

func (service *TaskService) GetTaskByID(param dtos.GetAnyFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()

	return service.repo.GetByAny(contextServer, param)
}

func (s *TaskService) GetPagination(params dtosPage.PaginationParams) dtosPage.PaginationResult[entities.Task] {
	contextServer := utils.CreateContextServerWithTimeout()
	resultContext := s.repo.GetPagination(contextServer, params)

	return dtosPage.PaginationResult[entities.Task]{
		Items: *params.Result.(*[]entities.Task),
		PageCurrent: params.Skip,
		PaginationResultContext: resultContext,
	}
}

func (service *TaskService) UpdateTask(id, userID primitive.ObjectID, dto dtos.Task) error {
	contextServer := utils.CreateContextServerWithTimeout()
	task := entities.Task{
		Title: dto.Title,
		Description: dto.Description,
		UpdateAt: utils.TimeNowBrazil(),
	}

	params := dtos.UpdateFilter{
		ID: id,
		Dto: task,
		ForeignKey: "userID",
		ForeignKeyValue: userID,
	}

	err := service.repo.Update(contextServer, params)
	if err != nil {
		return err
	}

	return nil
}

func (service *TaskService) DeleteTask(params dtos.DeleteFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()
	
	err := service.repo.Delete(contextServer, params)
	if err != nil {
		return err
	}

	return nil
}