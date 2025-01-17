package repository

import (
	"1uciuszzz/todo-app-golang/internal/model"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *model.Todo) error
	FindByID(id uint) (*model.Todo, error)
	FindAll() ([]model.Todo, error)
	Update(todo *model.Todo) error
	Delete(id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) FindByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) FindAll() ([]model.Todo, error) {
	var todos []model.Todo
	if err := r.db.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) Update(todo *model.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&model.Todo{}, id).Error
}
