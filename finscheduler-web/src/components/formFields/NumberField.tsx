import {Field, NumberInput} from "@chakra-ui/react";

interface NumberFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    defaultValue?: string;
    step?: number;
    min?: number;
    max?: number;
    required?: boolean;
}

export default function NumberField({
    label,
    value,
    onChange,
    defaultValue,
    step,
    min,
    max,
    required = false,
}: NumberFieldProps) {
    return (
        <Field.Root required={required} gap={0}>
            <Field.Label color="neon.blue">{label}</Field.Label>
            <NumberInput.Root
                value={value}
                onValueChange={(details) => onChange(details.value)}
                bg="bg.layer2"
                borderColor="glass.border"
                color="neon.blue"
                defaultValue={defaultValue}
                step={step}
                min={min}
                max={max}
                width="100%"
                gap={0}
            >
                <NumberInput.Control>
                    <NumberInput.IncrementTrigger
                        bg="bg.layer2"
                        color="neon.blue"
                        p={0}
                        _hover={{
                            filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))",
                            boxShadow: "0 0 12px rgba(0, 212, 255, 0.45)",
                        }}
                    >
                        +
                    </NumberInput.IncrementTrigger>
                    <NumberInput.DecrementTrigger
                        bg="bg.layer2"
                        color="neon.blue"
                        p={0}
                        _hover={{
                            filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))",
                            boxShadow: "0 0 12px rgba(0, 212, 255, 0.45)",
                        }}
                    >
                        -
                    </NumberInput.DecrementTrigger>
                </NumberInput.Control>
                <NumberInput.Input _placeholder={{color: 'text.placeholder'}}/>
            </NumberInput.Root>
        </Field.Root>
    );
}
