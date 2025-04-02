// components/AuthForm.jsx
import { useState } from 'react';
import { useLocation } from 'react-router-dom';
import StyledButton from '../ui/styled-button';
import StyledInput from '../ui/styled-input';

const AuthForm = () => {
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const authChallenge = queryParams.get('login_challenge');
    const error = queryParams.get('error');

    const [login, setLogin] = useState('');
    const [pass, setPass] = useState('');
    const [notAvailable, setNotAvailable] = useState(false);

    const handleSubmit = async (e) => {
        setNotAvailable(true);
        setTimeout(() => {
            setNotAvailable(false);
        }, 3000);
    };

    return (
        <div className="center">
            {error && <p>{error}</p>}

            <form onSubmit={handleSubmit} action="/api/auth/login" method="POST" className="form">
                <div className="title">
                    Введите логин и пароль
                </div>

                <input type="hidden" name="challenge" value={authChallenge} />

                <StyledInput
                    id="login"
                    value={login}
                    onChange={(e) => setLogin(e.target.value)}
                    type="text"
                    placeholder="Логин"
                    name="login"
                />

                <StyledInput
                    id="password"
                    value={pass}
                    onChange={(e) => setPass(e.target.value)}
                    type="password"
                    placeholder="Пароль"
                    name="password"
                />

                <StyledButton type="submit" disabled={notAvailable}>
                    Далее
                </StyledButton>
            </form>
        </div>
    );
};

export default AuthForm;