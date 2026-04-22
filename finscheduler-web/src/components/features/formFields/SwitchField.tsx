import {Field, Flex, Switch} from "@chakra-ui/react";

interface SwitchFieldProps {
    label: string;
    checked: boolean;
    onChange: (checked: boolean) => void;
}

export default function SwitchField({label, checked, onChange}: SwitchFieldProps) {
    return (
        <Field.Root gap={0}>
            <Flex align="center" gap={2}>
                <Field.Label color="neon.blue" mb={0}>
                    {label}
                </Field.Label>
                <Switch.Root
                    checked={checked}
                    onCheckedChange={(details) => onChange(details.checked)}
                >
                    <Switch.HiddenInput/>
                    <Switch.Control
                        bg={checked ? "neon.blue" : "neon.purple"}
                        filter={checked
                            ? "drop-shadow(0 0 8px rgba(0, 212, 255, 0.9))"
                            : "drop-shadow(0 0 8px rgba(212, 0, 255, 0.9))"}
                        boxShadow={checked
                            ? "0 0 12px rgba(0, 212, 255, 0.6)"
                            : "0 0 12px rgba(212, 0, 255, 0.6)"}
                        transition="all 0.3s ease-in-out"
                    >
                        <Switch.Thumb/>
                    </Switch.Control>
                </Switch.Root>
            </Flex>
        </Field.Root>
    );
}
