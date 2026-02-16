import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Login from './pages/Login';
import Main from './pages/Main';
import Products from './pages/Products';
import Stores from './pages/Stores';
import Reports from './pages/Reports';

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/main" element={<Main />} />
                    <Route path="/products" element={<Products />} />
                    <Route path="/stores" element={<Stores />} />
                    <Route path="/reports" element={<Reports />} />
                    <Route path="/" element={<Navigate to="/main" replace />} />
                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;
