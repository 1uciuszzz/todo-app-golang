package service

import (
	"1uciuszzz/todo-app-golang/internal/model"
	"1uciuszzz/todo-app-golang/internal/repository"
)

type TodoService interface {
	CreateTodo(todo *model.Todo) error
	GetTodo(id uint) (*model.Todo, error)
	ListTodos() ([]model.Todo, error)
	UpdateTodo(todo *model.Todo) error
	DeleteTodo(id uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(todo *model.Todo) error {
	return s.repo.Create(todo)
}

func (s *todoService) GetTodo(id uint) (*model.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) ListTodos() ([]model.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) UpdateTodo(todo *model.Todo) error {
	return s.repo.Update(todo)
}

func (s *todoService) DeleteTodo(id uint) error {
	return s.repo.Delete(id)
}
