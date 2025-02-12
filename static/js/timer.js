class TaskTimer {
    constructor(taskId) {
        this.taskId = taskId;
        this.startTime = null;
        this.elapsedTime = 0;
        this.timerId = null;
        this.isRunning = false;

        const timerElement = document.getElementById(`timer-${taskId}`);
        if (timerElement) {
            const initialTime = parseInt(timerElement.dataset.initialTime || '0', 10);
            this.elapsedTime = initialTime * 1000; // Конвертируем в миллисекунды
        }
    }

    start() {
        if (!this.isRunning) {
            this.startTime = Date.now() - this.elapsedTime;
            this.isRunning = true;
            this.updateDisplay();
            this.timerId = setInterval(() => this.updateDisplay(), 1000);
        }
    }

    pause() {
        if (this.isRunning) {
            clearInterval(this.timerId);
            this.timerId = null;
            this.elapsedTime = Date.now() - this.startTime;
            this.isRunning = false;

            if (this.elapsedTime > 0) {
                this.saveTimeToServer(Math.floor(this.elapsedTime / 1000));
            }
        }
    }

    updateDisplay() {
        const currentTime = Date.now();
        this.elapsedTime = currentTime - this.startTime;
        const seconds = Math.floor(this.elapsedTime / 1000);
        const timerElement = document.getElementById(`timer-${this.taskId}`);
        if (timerElement) {
            timerElement.textContent = this.formatTime(seconds);
        }
    }

    formatTime(totalSeconds) {
        if (isNaN(totalSeconds)) return "00:00:00";

        const hours = Math.floor(totalSeconds / 3600);
        const minutes = Math.floor((totalSeconds % 3600) / 60);
        const seconds = totalSeconds % 60;

        return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
    }

    async saveTimeToServer(timeInSeconds) {
        if (isNaN(timeInSeconds)) {
            console.error('Invalid time value:', timeInSeconds);
            return;
        }

        try {
            const response = await fetch(`/tasks/${this.taskId}/time`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `time=${timeInSeconds}`
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Failed to save time:', errorText);
                alert('Failed to save time: ' + errorText);
                return;
            }

            const data = await response.json();

            // Обновляем значение времени на основе ответа сервера
            if (data.timeSpent !== undefined) {
                this.elapsedTime = data.timeSpent * 1000; // Конвертируем в миллисекунды
                this.updateDisplay();
            }
        } catch (error) {
            console.error('Error saving time:', error);
            alert('Error saving time: ' + error.message);
        }
    }
}

// Инициализация и управление таймерами
const timers = new Map();

function initializeTimer(taskId) {
    if (timers.has(taskId)) return;

    const timerElement = document.getElementById(`timer-${taskId}`);
    if (!timerElement) return;

    const initialTime = parseInt(timerElement.dataset.initialTime || '0', 10);
    const status = timerElement.dataset.status;

    const timer = new TaskTimer(taskId);
    timer.elapsedTime = initialTime * 1000;
    timers.set(taskId, timer);

    // Автозапуск таймера, если задача в статусе STARTED
    if (status === 'STARTED') {
        timer.start();
        const button = document.querySelector(`button[data-task-id="${taskId}"]`);
        if (button) {
            button.dataset.timerState = 'running';
        }
    }

    return timer;
}

// Инициализация всех таймеров при загрузке страницы
function initializeAllTimers() {
    document.querySelectorAll('[data-task-id]').forEach(element => {
        const taskId = element.dataset.taskId;
        initializeTimer(taskId);
    });
}

// Обработка кликов по кнопкам "Play/Pause"
function handleTimerButtonClick(e) {
    const button = e.target.closest('[data-action="toggle-timer"]');
    if (!button) return;

    const taskId = button.dataset.taskId;
    let timer = timers.get(taskId);

    if (!timer) {
        timer = initializeTimer(taskId);
    }

    const timerState = button.dataset.timerState;

    if (timerState === 'stopped') {
        timer.start();
        button.dataset.timerState = 'running';
    } else {
        timer.pause();
        button.dataset.timerState = 'stopped';
    }
}

// Обработка событий HTMX
function handleHtmxBeforeSwap(evt) {
    const taskId = evt.detail.target.id.replace('task-', '');
    const timer = timers.get(taskId);
    if (timer && timer.isRunning) {
        const timerElement = document.getElementById(`timer-${taskId}`);
        if (timerElement) {
            timerElement.dataset.elapsedTime = Math.floor(timer.elapsedTime / 1000);
        }
    }
}

function handleHtmxAfterSwap(evt) {
    const taskId = evt.detail.target.id.replace('task-', '');
    const task = document.getElementById(`task-${taskId}`);
    if (!task) return;

    const status = task.querySelector('[data-status]')?.dataset.status;
    
    if (status === 'COMPLETED') {
        const timer = timers.get(taskId);
        if (timer && timer.isRunning) {
            timer.pause();
        }
        // Блокируем кнопку Play
        const playButton = task.querySelector('[data-action="toggle-timer"]');
        if (playButton) {
            playButton.disabled = true;
            playButton.classList.add('opacity-50', 'cursor-not-allowed');
        }
    } else {
        // инициализируем таймер с актуальным значением времени
        const timer = timers.get(taskId);
        if (!timer || !timer.isRunning) {
            initializeTimer(taskId);
        }
    }
}

// Инициализация всех обработчиков событий
function initializeEventHandlers() {
    window.addEventListener('load', initializeAllTimers);
    document.addEventListener('click', handleTimerButtonClick);
    document.addEventListener('htmx:beforeSwap', handleHtmxBeforeSwap);
    document.addEventListener('htmx:afterSwap', handleHtmxAfterSwap);
}

// Запуск инициализации
initializeEventHandlers();