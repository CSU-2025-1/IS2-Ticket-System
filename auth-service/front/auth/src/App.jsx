// App.jsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AuthForm from './components/auth';
import ConsentForm from './components/consent';

function App() {
  return (
      <div className="wrapper">
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<AuthForm />} />
            <Route path="/consent" element={<ConsentForm />} />
          </Routes>
        </BrowserRouter>
      </div>
  );
}

export default App;