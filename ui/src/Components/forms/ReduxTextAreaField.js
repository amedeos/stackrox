import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'redux-form';

const ReduxTextAreaField = ({ name, disabled, placeholder }) => (
    <Field
        key={name}
        name={name}
        component="textarea"
        className="border rounded-l p-3 border-base-300 text-base-600 w-full font-400"
        disabled={disabled}
        rows={4}
        placeholder={placeholder}
    />
);

ReduxTextAreaField.propTypes = {
    name: PropTypes.string.isRequired,
    disabled: PropTypes.bool,
    placeholder: PropTypes.string.isRequired
};

ReduxTextAreaField.defaultProps = {
    disabled: false
};

export default ReduxTextAreaField;
