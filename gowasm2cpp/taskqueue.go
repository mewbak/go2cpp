// SPDX-License-Identifier: Apache-2.0

package gowasm2cpp

import (
	"os"
	"path/filepath"
	"text/template"
)

func writeTaskQueue(dir string, incpath string, namespace string) error {
	{
		f, err := os.Create(filepath.Join(dir, "taskqueue.h"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := taskqueueHTmpl.Execute(f, struct {
			IncludeGuard string
			IncludePath  string
			Namespace    string
		}{
			IncludeGuard: includeGuard(namespace) + "_TASKQUEUE_H",
			IncludePath:  incpath,
			Namespace:    namespace,
		}); err != nil {
			return err
		}
	}
	{
		f, err := os.Create(filepath.Join(dir, "taskqueue.cpp"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := taskqueueCppTmpl.Execute(f, struct {
			IncludePath string
			Namespace   string
		}{
			IncludePath: incpath,
			Namespace:   namespace,
		}); err != nil {
			return err
		}
	}
	return nil
}

var taskqueueHTmpl = template.Must(template.New("taskqueue.h").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#ifndef {{.IncludeGuard}}
#define {{.IncludeGuard}}

#include <mutex>
#include <condition_variable>
#include <functional>
#include <queue>
#include <thread>

namespace {{.Namespace}} {

class TaskQueue {
public:
  using Task = std::function<void()>;

  void Enqueue(Task task);
  Task Dequeue();

private:
  std::mutex mutex_;
  std::condition_variable cond_;
  std::queue<Task> queue_;
};

class Timer {
public:
  Timer(std::function<void()> func, double interval);
  ~Timer();

  void Stop();

private:
  enum class Result {
    Timeout,
    NoTimeout,
  };

  Result WaitFor(double milliseconds);

  // A mutex and a condition must be constructed before the thread starts.
  std::mutex mutex_;
  std::condition_variable cond_;
  std::thread thread_;
  bool stopped_ = false;
};

}

#endif  // {{.IncludeGuard}}
`))

var taskqueueCppTmpl = template.Must(template.New("taskqueue.cpp").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "{{.IncludePath}}taskqueue.h"

#include <chrono>
#include <memory>

namespace {{.Namespace}} {

void TaskQueue::Enqueue(Task task) {
  {
    std::lock_guard<std::mutex> lock{mutex_};
    queue_.push(task);
  }
  cond_.notify_one();
}

TaskQueue::Task TaskQueue::Dequeue() {
  std::unique_lock<std::mutex> lock{mutex_};
  cond_.wait(lock, [this]{ return !queue_.empty(); });
  Task task = queue_.front();
  queue_.pop();
  return task;
}

Timer::Timer(std::function<void()> func, double interval)
    : thread_{[this, interval](std::function<void()> func) {
        Result result = WaitFor(interval);
        if (result == Timer::Result::NoTimeout) {
          return;
        }
        func();
      }, std::move(func)} {
}

Timer::~Timer() {
  if (thread_.joinable()) {
    thread_.join();
  }
}

void Timer::Stop() {
  {
    std::lock_guard<std::mutex> lock{mutex_};
    stopped_ = true;
  }
  cond_.notify_one();
}

Timer::Result Timer::WaitFor(double milliseconds) {
  std::unique_lock<std::mutex> lock{mutex_};
  auto duration = std::chrono::duration<double, std::milli>(milliseconds);
  bool result = cond_.wait_for(lock, duration, [this]{ return stopped_; });
  return result ? Result::NoTimeout : Result::Timeout;
}

}
`))
