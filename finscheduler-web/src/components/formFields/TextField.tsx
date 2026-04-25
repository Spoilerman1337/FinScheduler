import {Field, Input} from "@chakra-ui/react";

interface TextFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    required?: boolean;
}

export default function TextField({
    label,
    value,
    onChange,
    placeholder,
    required = false,
}: TextFieldProps) {
    return (
        <Field.Root required={required} gap={0}>
            <Field.Label color="neon.blue">{label}</Field.Label>
            <Input
                value={value}
                onChange={(e) => onChange(e.target.value)}
                bg="bg.layer2"
                borderColor="glass.border"
                color="neon.blue"
                _placeholder={{color: 'text.placeholder'}}
                placeholder={placeholder}
            />
        </Field.Root>
    );
}
