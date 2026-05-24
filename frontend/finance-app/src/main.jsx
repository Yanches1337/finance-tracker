import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
// Проверьте путь! Если папка context лежит в src, то путь именно такой:
import { AuthProvider } from './context/AuthContext.jsx'

ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        <AuthProvider>
            <App />
        </AuthProvider>
    </React.StrictMode>,
)