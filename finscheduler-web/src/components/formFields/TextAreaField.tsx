import {Field, Textarea} from '@chakra-ui/react';

interface TextareaFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    rows?: number;
    required?: boolean;
    invalid?: boolean;
    errorText?: string;
}

export default function TextAreaField({
    label,
    value,
    onChange,
    placeholder,
    rows = 4,
    required = false,
    invalid = false,
    errorText,
}: TextareaFieldProps) {
    const isInvalid = invalid || Boolean(errorText);

    return (
        <Field.Root required={required} invalid={isInvalid} gap={1}>
            <Field.Label color="neon.blue">{label}</Field.Label>
            <Textarea
                value={value}
                onChange={(e) => onChange(e.target.value)}
                bg="bg.layer2"
                borderColor={isInvalid ? 'border.error' : 'glass.border'}
                color="neon.blue"
                _placeholder={{color: 'text.placeholder'}}
                placeholder={placeholder}
                rows={rows}
                _focusVisible={{
                    borderColor: isInvalid ? 'border.error' : 'neon.blue',
                    boxShadow: isInvalid
                        ? '0 0 0 1px rgba(255, 74, 122, 0.45)'
                        : '0 0 0 1px rgba(0, 212, 255, 0.35)',
                }}
            />
            {errorText ? <Field.ErrorText color="fg.error">{errorText}</Field.ErrorText> : null}
        </Field.Root>
    );
}
