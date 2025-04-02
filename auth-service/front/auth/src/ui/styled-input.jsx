// ui/StyledInput.jsx
const StyledInput = ({ id, value, onChange, type, placeholder, name }) => {
    return (
        <input
            id={id}
            value={value}
            onChange={onChange}
            type={type}
            placeholder={placeholder}
            name={name}
        />
    );
};

export default StyledInput;