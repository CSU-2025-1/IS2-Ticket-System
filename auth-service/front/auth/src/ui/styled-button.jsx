// ui/StyledButton.jsx
const StyledButton = ({ children, disabled, type }) => {
    return (
        <button type={type} disabled={disabled}>
            {children}
        </button>
    );
};

export default StyledButton;