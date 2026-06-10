import {Button, CloseButton, Dialog, Portal, Stack, Text} from '@chakra-ui/react';

interface UnsavedChangesDialogProps {
    open: boolean;
    onStay: () => void;
    onLeave: () => void;
}

export default function UnsavedChangesDialog({
    open,
    onStay,
    onLeave,
}: UnsavedChangesDialogProps) {
    return (
        <Dialog.Root
            open={open}
            onOpenChange={(details) => {
                if (!details.open) {
                    onStay();
                }
            }}
            placement="center"
            lazyMount
            unmountOnExit
        >
            <Portal>
                <Dialog.Backdrop />
                <Dialog.Positioner>
                    <Dialog.Content
                        bg="bg.layer1"
                        border="1px solid"
                        borderColor="glass.border"
                        boxShadow="0 0 24px rgba(255, 166, 0, 0.14)"
                        maxW="560px"
                    >
                        <Dialog.Header>
                            <Dialog.Title color="fg">Есть несохранённые изменения</Dialog.Title>
                            <Dialog.CloseTrigger asChild>
                                <CloseButton
                                    color="fg"
                                    _hover={{bg: 'rgba(255, 255, 255, 0.08)'}}
                                />
                            </Dialog.CloseTrigger>
                        </Dialog.Header>
                        <Dialog.Body>
                            <Stack gap={3}>
                                <Text color="fg">
                                    Если уйти со страницы сейчас, несохранённые изменения будут
                                    потеряны.
                                </Text>
                                <Text color="fg.muted">
                                    Сохраните изменения или подтвердите уход без сохранения.
                                </Text>
                            </Stack>
                        </Dialog.Body>
                        <Dialog.Footer>
                            <Button variant="ghost" onClick={onStay}>
                                Остаться
                            </Button>
                            <Button
                                bg="orange.400"
                                color="black"
                                _hover={{bg: 'orange.300'}}
                                onClick={onLeave}
                            >
                                Закрыть без сохранения
                            </Button>
                        </Dialog.Footer>
                    </Dialog.Content>
                </Dialog.Positioner>
            </Portal>
        </Dialog.Root>
    );
}
