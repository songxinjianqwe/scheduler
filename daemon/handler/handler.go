package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/songxinjianqwe/scheduler/common"
	"github.com/songxinjianqwe/scheduler/daemon/engine"
	"github.com/songxinjianqwe/scheduler/daemon/engine/standalone"
	"io/ioutil"
	"net/http"
	"strconv"
)

var scheduler engine.Engine

func init() {
	// 暂时使用这个实现类8
	scheduler = standalone.NewStandAloneEngine()
}

func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tasks, err := scheduler.List()
	if err != nil {
		http.Error(w, "获取任务列表失败", http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "获取任务列表失败", http.StatusInternalServerError)
		return
	}
	//WriteHeader只能调用一次，否则会引起http: multiple response.WriteHeader calls
	//WriteHeader必须在Write()之前调用，因为在Write()调用过程中，如果发现WriteHeader没有调用过，那么它会自行写入一次200作为默认header
	//对于标准Header如ContentType，Header().Set()必须在WriteHeader/Write之前调用，否则不会生效。因为按照写入顺序，是 Header -> StatusHeader -> Body -> Trailer Header.
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func GetTaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	watch, err := strconv.ParseBool(r.URL.Query().Get("watch"))
	if err != nil {
		http.Error(w, fmt.Sprintf("watch参数必须传入，且为true或false: %s", err.Error()), http.StatusBadRequest)
		return
	}
	versionStr := r.URL.Query().Get("version")
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("version参数必须传入，且为int64类型: %s", err.Error()), http.StatusBadRequest)
		return
	}
	task, err := scheduler.Get(id, watch, version)
	if err != nil {
		http.Error(w, fmt.Sprintf("任务ID[%s]不存在", id), http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "序列化任务失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func SubmitTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "读取任务Body失败", http.StatusBadRequest)
		return
	}
	var task common.Task
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		http.Error(w, "反序列任务失败", http.StatusInternalServerError)
		return
	}
	log.Infof("Receive a task:%#v", task)
	err = scheduler.Submit(task)
	if err != nil {
		log.Errorf("提交任务失败，失败原因: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func StopTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	err := scheduler.Stop(id)
	if err != nil {
		log.Errorf("停止任务失败，失败原因: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	err := scheduler.Delete(id)
	if err != nil {
		log.Errorf("删除任务失败，失败原因: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
