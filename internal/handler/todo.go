package handler

import (
	"net/http"
	"strconv"

	"1uciuszzz/todo-app-golang/internal/model"
	"1uciuszzz/todo-app-golang/internal/service"
	"1uciuszzz/todo-app-golang/pkg/response"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) Create(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CreateTodo(&todo); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, todo)
}

func (h *TodoHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	todo, err := h.service.GetTodo(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Todo not found")
		return
	}

	response.Success(c, todo)
}

func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.service.ListTodos()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, todos)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	todo.ID = uint(id)
	if err := h.service.UpdateTodo(&todo); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.DeleteTodo(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Todo deleted successfully"})
}
