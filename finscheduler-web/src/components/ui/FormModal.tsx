import {
    Button,
    CloseButton,
    Dialog,
    Portal,
    Stack,
    Text,
} from "@chakra-ui/react";
import type {ReactNode} from "react";

interface FormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: () => void | Promise<void>;
    title: string;
    children: ReactNode;
    error?: string | null;
    loading?: boolean;
}

const MODAL_MAX_WIDTH = "600px";
const MODAL_CONTENT_GAP = 4;
const SUBMIT_TEXT = "Сохранить";
const CANCEL_TEXT = "Отмена";

export default function FormModal({
    isOpen,
    onClose,
    onSubmit,
    title,
    children,
    error = null,
    loading = false,
}: FormModalProps) {
    return (
        <Dialog.Root open={isOpen} onOpenChange={(details) => !details.open && onClose()} placement="center">
            <Portal>
                <Dialog.Backdrop/>
                <Dialog.Positioner>
                    <Dialog.Content bg="bg.layer1" border="1px solid" borderColor="glass.border" maxW={MODAL_MAX_WIDTH}>
                        <Dialog.Header>
                            <Dialog.Title color="neon.blue">
                                {title}
                            </Dialog.Title>
                            <Dialog.CloseTrigger asChild bg="bg.layer1" border="1px solid" borderColor="neon.blue">
                                <CloseButton
                                    color="neon.blue"
                                    bg="transparent"
                                    border="none"
                                    outline="none"
                                    boxShadow="none"
                                    transition="all 0.3s ease-in-out"
                                    _hover={{
                                        color: "neon.blue",
                                        bg: "glass.bgHover",
                                        backdropFilter: "blur(12px)",
                                        borderColor: "neon.blue",
                                        outline: "none",
                                        filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))",
                                        boxShadow: "0 0 12px rgba(0, 212, 255, 0.45)",
                                    }}
                                    _focusVisible={{
                                        bg: "transparent",
                                        outline: "none",
                                        boxShadow: "none",
                                    }}
                                />
                            </Dialog.CloseTrigger>
                        </Dialog.Header>

                        <Dialog.Body>
                            {error && (
                                <Text color="neon.pink" fontSize="sm" mb={4}>
                                    {error}
                                </Text>
                            )}
                            <Stack gap={MODAL_CONTENT_GAP}>
                                {children}
                            </Stack>
                        </Dialog.Body>

                        <Dialog.Footer>
                            <Button
                                variant="ghost"
                                mr={3}
                                onClick={onClose}
                                color="textMuted"
                                _hover={{bg: 'bg.layer2'}}
                            >
                                {CANCEL_TEXT}
                            </Button>
                            <Button
                                bg="neon.blue"
                                color="bg.base"
                                onClick={onSubmit}
                                loading={loading}
                                _hover={{bg: 'neon.blue', opacity: 0.8}}
                            >
                                {SUBMIT_TEXT}
                            </Button>
                        </Dialog.Footer>
                    </Dialog.Content>
                </Dialog.Positioner>
            </Portal>
        </Dialog.Root>
    );
}
