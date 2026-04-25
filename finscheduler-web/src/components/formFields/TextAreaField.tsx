import {Field, Textarea} from "@chakra-ui/react";

interface TextareaFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    rows?: number;
    required?: boolean;
}

export default function TextAreaField({
    label,
    value,
    onChange,
    placeholder,
    rows = 4,
    required = false,
}: TextareaFieldProps) {
    return (
        <Field.Root required={required} gap={0}>
            <Field.Label color="neon.blue">{label}</Field.Label>
            <Textarea
                value={value}
                onChange={(e) => onChange(e.target.value)}
                bg="bg.layer2"
                borderColor="glass.border"
                color="neon.blue"
                _placeholder={{color: 'text.placeholder'}}
                placeholder={placeholder}
                rows={rows}
            />
        </Field.Root>
    );
}
