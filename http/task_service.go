package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"

	"github.com/influxdata/platform"
	kerrors "github.com/influxdata/platform/kit/errors"
	"github.com/julienschmidt/httprouter"
)

const taskPath = "/v1/tasks"

// TaskHandler represents an HTTP API handler for tasks.
type TaskHandler struct {
	*httprouter.Router
	TaskService platform.TaskService
}

// NewTaskHandler returns a new instance of TaskHandler.
func NewTaskHandler() *TaskHandler {
	h := &TaskHandler{
		Router: httprouter.New(),
	}

	h.HandlerFunc("GET", "/v1/tasks", h.handleGetTasks)
	h.HandlerFunc("POST", "/v1/tasks", h.handlePostTask)

	h.HandlerFunc("GET", "/v1/tasks/:tid", h.handleGetTask)
	h.HandlerFunc("PATCH", "/v1/tasks/:tid", h.handleUpdateTask)
	h.HandlerFunc("DELETE", "/v1/tasks/:tid", h.handleDeleteTask)

	h.HandlerFunc("GET")

	h.HandlerFunc("GET", "/v1/tasks/:tid/logs", h.handleGetLogs)
	h.HandlerFunc("GET", "/v1/tasks/:tid/runs/:rid/logs", h.handleGetLogs)

	h.HandlerFunc("GET", "/v1/tasks/:tid/runs", h.handleGetRuns)
	h.HandlerFunc("GET", "/v1/tasks/:tid/runs/:rid", h.handleGetRun)
	h.HandlerFunc("POST", "/v1/tasks/:tid/runs/:rid/retry", h.handleRetryRun)

	return h
}

func (h *TaskHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeGetTasksRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	tasks, _, err := h.TaskService.FindTasks(ctx, req.filter)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, tasks); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type getTasksRequest struct {
	filter platform.TaskFilter
}

func decodeGetTasksRequest(ctx context.Context, r *http.Request) (*getTasksRequest, error) {
	qp := r.URL.Query()
	req := &getTasksRequest{}

	if id := qp.Get("after"); id != "" {
		req.filter.After = &platform.ID{}
		if err := req.filter.After.DecodeFromString(id); err != nil {
			return nil, err
		}
	}

	if id := qp.Get("organization"); id != "" {
		req.filter.Organization = &platform.ID{}
		if err := req.filter.Organization.DecodeFromString(id); err != nil {
			return nil, err
		}
	}

	if id := qp.Get("user"); id != "" {
		req.filter.User = &platform.ID{}
		if err := req.filter.User.DecodeFromString(id); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (h *TaskHandler) handlePostTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodePostTaskRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := h.TaskService.CreateTask(ctx, req.Task); err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusCreated, req.Task); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type postTaskRequest struct {
	Task *platform.Task
}

func decodePostTaskRequest(ctx context.Context, r *http.Request) (*postTaskRequest, error) {
	task := &platform.Task{}
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		return nil, err
	}
	return &postTaskRequest{
		Task: task,
	}, nil
}

func (h *TaskHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeGetTaskRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	task, err := h.TaskService.FindTaskByID(ctx, req.TaskID)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, task); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type getTaskRequest struct {
	TaskID platform.ID
}

func decodeGetTaskRequest(ctx context.Context, r *http.Request) (*getTaskRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("tid")
	if id == "" {
		return nil, kerrors.InvalidDataf("url missing id")
	}

	var i platform.ID
	if err := i.DecodeFromString(id); err != nil {
		return nil, err
	}

	req := &getTaskRequest{
		TaskID: i,
	}

	return req, nil
}

func (h *TaskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeUpdateTaskRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	task, err := h.TaskService.UpdateTask(ctx, req.TaskID, req.Update)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, task); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type updateTaskRequest struct {
	Update platform.TaskUpdate
	TaskID platform.ID
}

func decodeUpdateTaskRequest(ctx context.Context, r *http.Request) (*updateTaskRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("tid")
	if id == "" {
		return nil, kerrors.InvalidDataf("you must provide a task ID")
	}

	var i platform.ID
	if err := i.DecodeFromString(id); err != nil {
		return nil, err
	}

	var upd platform.TaskUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		return nil, err
	}

	return &updateTaskRequest{
		Update: upd,
		TaskID: i,
	}, nil
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeDeleteTaskRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := h.TaskService.DeleteTask(ctx, req.TaskID); err != nil {
		EncodeError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type deleteTaskRequest struct {
	TaskID platform.ID
}

func decodeDeleteTaskRequest(ctx context.Context, r *http.Request) (*deleteTaskRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("tid")
	if id == "" {
		return nil, kerrors.InvalidDataf("you must provide a task ID")
	}

	var i platform.ID
	if err := i.DecodeFromString(id); err != nil {
		return nil, err
	}

	return &deleteTaskRequest{
		TaskID: i,
	}, nil
}

func (h *TaskHandler) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeGetLogsRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	logs, _, err := h.TaskService.FindLogs(ctx, req.filter)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, logs); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type getLogsRequest struct {
	filter platform.LogFilter
}

func decodeGetLogsRequest(ctx context.Context, r *http.Request) (*getLogsRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("tid")
	if id == "" {
		return nil, kerrors.InvalidDataf("you must provide a task ID")
	}

	req := &getLogsRequest{}
	req.filter.Task = &platform.ID{}
	if err := req.filter.Task.DecodeFromString(id); err != nil {
		return nil, err
	}

	if id := params.ByName("rid"); id != "" {
		req.filter.Run = &platform.ID{}
		if err := req.filter.Run.DecodeFromString(id); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (h *TaskHandler) handleGetRuns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeGetRunsRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	runs, _, err := h.TaskService.FindRuns(ctx, req.filter)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, runs); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type getRunsRequest struct {
	filter platform.RunFilter
}

func decodeGetRunsRequest(ctx context.Context, r *http.Request) (*getRunsRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("tid")
	if id == "" {
		return nil, kerrors.InvalidDataf("you must provide a task ID")
	}

	req := &getRunsRequest{}
	req.filter.Task = &platform.ID{}
	if err := req.filter.Task.DecodeFromString(id); err != nil {
		return nil, err
	}

	qp := r.URL.Query()

	if id := qp.Get("after"); id != "" {
		req.filter.After = &platform.ID{}
		if err := req.filter.After.DecodeFromString(id); err != nil {
			return nil, err
		}
	}

	if limit := qp.Get("limit"); limit != "" {
		i, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}

		if i < 1 || i > 100 {
			return nil, kerrors.InvalidDataf("limit must be between 1 and 100")
		}

		req.filter.Limit = i
	}

	if time := qp.Get("afterTime"); time != "" {
		// TODO (jm): verify valid RFC3339
		req.filter.AfterTime = time
	}

	if time := qp.Get("beforeTime"); time != "" {
		// TODO (jm): verify valid RFC3339
		req.filter.BeforeTime = time
	}

	return req, nil
}

func (h *TaskHandler) handleGetRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeGetRunRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	run, err := h.TaskService.FindRunByID(ctx, req.RunID)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, run); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type getRunRequest struct {
	RunID platform.ID
}

func decodeGetRunRequest(ctx context.Context, r *http.Request) (*getRunRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("rid")
	if id == "" {
		return nil, kerrors.InvalidDataf("you must provide a run ID")
	}

	var i platform.ID
	if err := i.DecodeFromString(id); err != nil {
		return nil, err
	}

	return &getRunRequest{
		RunID: i,
	}, nil
}

func (h *TaskHandler) handleRetryRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeRetryRunRequest(ctx, r)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	run, err := h.TaskService.RetryRun(ctx, req.RunID)
	if err != nil {
		EncodeError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, run); err != nil {
		EncodeError(ctx, err, w)
		return
	}
}

type retryRunRequest struct {
	RunID platform.ID
}

func decodeRetryRunRequest(ctx context.Context, r *http.Request) (*retryRunRequest, error) {
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("rid")
	if id == "" {
		return nil, kerrors.InvalidDataf("you must provide a run ID")
	}

	var i platform.ID
	if err := i.DecodeFromString(id); err != nil {
		return nil, err
	}

	return &retryRunRequest{
		RunID: i,
	}, nil
}

// TaskService connects to Influx via HTTP using tokens to manage tasks
type TaskService struct {
	Addr               string
	Token              string
	InsecureSkipVerify bool
}

func (s *TaskService) FindTaskByID(ctx context.Context, id platform.ID) (*platform.Task, error) {
	u, err := newURL(s.Addr, taskIDPath(id))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+s.Token)

	hc := newClient(u.Scheme, s.InsecureSkipVerify)
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	if err := CheckError(resp); err != nil {
		return nil, err
	}

	var task platform.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &task, nil
}

func (s *TaskService) FindTask(ctx context.Context, filter platform.TaskFilter) (*platform.Task, error) {
	tasks, n, err := s.FindTasks(ctx, filter)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, errors.New("found no matching tasks")
	}

	return tasks[0], nil
}

func (s *TaskService) FindTasks(ctx context.Context, filter platform.TaskFilter) ([]*platform.Task, int, error) {
	u, err := newURL(s.Addr, bucketPath)
	if err != nil {
		return nil, 0, err
	}

	query := u.Query()
	if filter.Organization != nil {
		query.Add("org", filter.Organization.String())
	}
	if filter.ID != nil {
		query.Add("id", filter.ID.String())
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", "Token "+s.Token)

	hc := newClient(u.Scheme, s.InsecureSkipVerify)
	resp, err := hc.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if err := CheckError(resp); err != nil {
		return nil, 0, err
	}

	var tasks []*platform.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	return tasks, len(tasks), nil
}

func (s *TaskService) CreateTask(ctx context.Context, task *platform.Task) error {
	u, err := newURL(s.Addr, taskPath)
	if err != nil {
		return err
	}

	octets, err := json.Marshal(task)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(octets))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+s.Token)

	hc := newClient(u.Scheme, s.InsecureSkipVerify)

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	// TODO(jsternberg): Should this check for a 201 explicitly?
	if err := CheckError(resp); err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(task); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id platform.ID, upd platform.TaskUpdate) (*platform.Task, error) {
	u, err := newURL(s.Addr, taskIDPath(id))
	if err != nil {
		return nil, err
	}

	octets, err := json.Marshal(upd)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", u.String(), bytes.NewReader(octets))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+s.Token)

	hc := newClient(u.Scheme, s.InsecureSkipVerify)

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	if err := CheckError(resp); err != nil {
		return nil, err
	}

	var task platform.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id platform.ID) error {
	u, err := newURL(s.Addr, taskIDPath(id))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+s.Token)

	hc := newClient(u.Scheme, s.InsecureSkipVerify)
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	return CheckError(resp)
}

func taskIDPath(id platform.ID) string {
	return path.Join(taskPath, id.String())
}