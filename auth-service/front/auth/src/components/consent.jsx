// components/ConsentForm.jsx
import { useState } from 'react';
import { useLocation } from 'react-router-dom';
import StyledButton from '../ui/styled-button';

const ConsentForm = () => {
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const consentChallenge = queryParams.get('consent_challenge');

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState(null);

    // URL вашего API endpoint
    const FORM_ACTION_URL = '/api/auth/consent';

    const handleSubmit = (e) => {
        // Не вызываем preventDefault, чтобы форма отправилась стандартным способом
        setIsSubmitting(true);
        setError(null);
    };

    return (
        <div className="center">
            <form
                action={FORM_ACTION_URL}
                method="POST"
                onSubmit={handleSubmit}
                className="form"
            >
                <div className="title">
                    Продолжая, вы соглашаетесь с условиями <br />
                    предоставления сервиса.
                </div>

                <input
                    type="hidden"
                    name="challenge"
                    value={consentChallenge || ''}
                />

                {error && (
                    <div className="error-message">
                        {error}
                    </div>
                )}

                <StyledButton
                    type="submit"
                    disabled={isSubmitting}
                >
                    {isSubmitting ? 'Отправка...' : 'Далее'}
                </StyledButton>
            </form>
        </div>
    );
};

export default ConsentForm;