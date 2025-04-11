document.addEventListener('DOMContentLoaded', () => {
    const eventsContainer = document.getElementById('events');
    const connectButton = document.getElementById('connect');
    const disconnectButton = document.getElementById('disconnect');
    let eventSource = null;

    function appendMessage(message, isError = false) {
        const div = document.createElement('div');
        div.textContent = message;
        div.className = isError ? 'error' : 'data';
        eventsContainer.appendChild(div);
        eventsContainer.scrollTop = eventsContainer.scrollHeight;
    }

    connectButton.addEventListener('click', () => {
        if (eventSource) {
            eventSource.close();
        }

        appendMessage('Connecting to SSE server...');
        eventSource = new EventSource('http://localhost:8080/event');
        
        // Enable disconnect button as soon as we attempt to connect
        connectButton.disabled = true;
        disconnectButton.disabled = false;

        eventSource.onopen = () => {
            appendMessage('Connection established');
        };

        eventSource.onmessage = (event) => {
            appendMessage(`Received data: ${event.data}`);
        };

        // This is equivalent to eventSource.onerror = (event) => {}
        // See https://developer.mozilla.org/en-US/docs/Web/API/EventSource/error_event
        eventSource.addEventListener('error', (event) => {
            // Check if the connection was closed normally (after 10 messages)
            if (eventSource.readyState === EventSource.CLOSED) {
                appendMessage('Connection closed by server (completed 10 messages)');
                resetButtons();
            } else if (eventSource.readyState === EventSource.CONNECTING) {
                // This happens when the server closes the connection but the browser tries to reconnect
                // We should close the connection explicitly to prevent continuous reconnection attempts
                appendMessage('Server closed the connection', true);
                eventSource.close();
                resetButtons();
            } else if (event.data) {
              appendMessage(`Received error: ${event.data}`, true);
          } else {
              appendMessage('An unknown error occurred', true);
          }
        });
    });

    disconnectButton.addEventListener('click', () => {
        if (eventSource) {
            eventSource.close();
            appendMessage('Connection closed by client');
            resetButtons();
        }
    });

    function resetButtons() {
        connectButton.disabled = false;
        disconnectButton.disabled = true;
    }
});
