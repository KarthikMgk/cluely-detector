// bg.js  (service-worker)

let port;

function connect() {

  // 1) Open native-messaging port
    port = chrome.runtime.connectNative("com.watchdog");
    
    port.onDisconnect.addListener(() => {
      console.warn("watchdog disconnected", chrome.runtime.lastError);
    });

  // 2) Log every message from the helper
    port.onMessage.addListener((msg) => {
        console.log("[watchdog]", msg);
        chrome.action.setBadgeText({ text: msg.overlay ? "⚠️" : "" });
    });

    // 3) Reconnect if the helper exits
    port.onDisconnect.addListener(() => setTimeout(connect, 1000));
}

// Kick things off as soon as the service-worker starts
connect();
