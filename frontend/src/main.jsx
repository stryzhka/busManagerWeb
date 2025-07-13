import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'

const container = document.getElementById('root')

const root = createRoot(container)

root.render(
    // <React.StrictMode>
        <App/>
    // </React.StrictMode>
)

document.addEventListener('click', (e) => {
    const target = e.target.closest('a');
    if (target && target.href) {
        e.preventDefault(); // Блокируем стандартное поведение
        window.runtime.BrowserOpenURL(target.href); // Открываем ссылку во внешнем браузере
    }
});