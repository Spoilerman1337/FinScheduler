import {Flex, Switch, Text} from "@chakra-ui/react";

interface SwitchFieldProps {
    label: string;
    checked: boolean;
    onChange: (checked: boolean) => void;
}

export default function SwitchField({label, checked, onChange}: SwitchFieldProps) {
    return (
        <Flex align="center" gap={2}>
            <Text color="neon.blue" mb={0} fontSize="sm" fontWeight="medium" lineHeight="1" userSelect="none">
                {label}
            </Text>
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
                    transition="background-color 0.3s ease-in-out, filter 0.3s ease-in-out, box-shadow 0.3s ease-in-out"
                    cursor="pointer"
                >
                    <Switch.Thumb
                        transitionProperty="translate"
                        transitionDuration="0.3s"
                        transitionTimingFunction="ease-in-out"
                    />
                </Switch.Control>
            </Switch.Root>
        </Flex>
    );
}
