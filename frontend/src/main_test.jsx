import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App_test from './App_test'

const container = document.getElementById('root')

const root = createRoot(container)

root.render(
    // <React.StrictMode>
    <App_test/>
    // </React.StrictMode>
)