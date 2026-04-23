import {
    createListCollection,
    Field,
    Flex,
    SelectContent,
    SelectHiddenSelect,
    SelectItem,
    SelectPositioner,
    SelectRoot,
    SelectTrigger,
    SelectValueText,
    Spinner,
} from "@chakra-ui/react";
import {ChevronDownIcon} from "lucide-react";
import {useMemo} from "react";
import type {SelectOption} from "./types.ts";

interface BaseSelectFieldProps {
    label: string;
    options: SelectOption[];
    placeholder?: string;
    loading?: boolean;
    required?: boolean;
}

type SingleSelectFieldProps = BaseSelectFieldProps & {
    multiple?: false;
    value: string;
    onChange: (value: string) => void;
};

type MultipleSelectFieldProps = BaseSelectFieldProps & {
    multiple: true;
    value: string[];
    onChange: (value: string[]) => void;
};

type SelectFieldProps = SingleSelectFieldProps | MultipleSelectFieldProps;

const SELECT_MAX_VISIBLE_OPTIONS = 5;
const SELECT_OPTION_HEIGHT_PX = 40;
const SELECT_CONTENT_MAX_HEIGHT = `${SELECT_MAX_VISIBLE_OPTIONS * SELECT_OPTION_HEIGHT_PX}px`;

export default function SelectField(props: SelectFieldProps) {
    const {
    label,
    options,
    placeholder,
    loading = false,
    required = false,
    } = props;
    const collection = useMemo(() => createListCollection({items: options}), [options]);
    const value = props.multiple ? props.value : props.value ? [props.value] : [];

    return (
        <Field.Root required={required} gap={0}>
            <Field.Label color="neon.blue">
                {label} {loading && <Spinner size="xs" ml={2}/>}
            </Field.Label>
            <SelectRoot
                multiple={props.multiple}
                collection={collection}
                value={value}
                onValueChange={(details) => {
                    if (props.multiple) {
                        props.onChange(details.value);
                    } else {
                        props.onChange(details.value[0] ?? '');
                    }
                }}
            >
                <SelectHiddenSelect/>
                <SelectTrigger asChild>
                    <Flex
                        align="center"
                        justify="space-between"
                        gap={2}
                        w="100%"
                        minH="40px"
                        px={3}
                        border="1px solid"
                        borderRadius="sm"
                        borderColor="glass.border"
                        bg="bg.layer2"
                        color="neon.blue"
                        transition="all 0.2s"
                        _hover={{
                            color: "neon.blue",
                            bg: "bg.layer2",
                            borderColor: "glass.border",
                            cursor: "pointer",
                        }}
                    >
                        <SelectValueText
                            placeholder={placeholder}
                            color="currentColor"
                            _placeholder={{color: 'textMuted'}}
                            overflow="hidden"
                            textOverflow="ellipsis"
                            whiteSpace="nowrap"
                            flex="1"
                        />
                        <ChevronDownIcon size={16} style={{flexShrink: 0}}/>
                    </Flex>
                </SelectTrigger>
                <SelectPositioner>
                    <SelectContent
                        p="1"
                        borderRadius="sm"
                        mt="1"
                        border="1px solid"
                        bg="bg.layer1"
                        borderColor="glass.borderStrong"
                        backdropFilter="blur(16px)"
                        boxShadow="lg"
                        zIndex="popover"
                        width="--trigger-width"
                        maxH={SELECT_CONTENT_MAX_HEIGHT}
                        overflowY="auto"
                        overscrollBehavior="contain"
                        className="custom-scrollbar"
                    >
                        {collection.items.map((option) => (
                            <SelectItem
                                item={option}
                                key={option.value}
                                display="flex"
                                alignItems="center"
                                py="1.5"
                                px="2"
                                borderRadius="sm"
                                color="neon.blue"
                                fontSize="sm"
                                transition="all 0.2s"
                                cursor="pointer"
                                _hover={{
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                    filter: "drop-shadow(0 0 8px rgba(212, 0, 255, 0.45))",
                                }}
                                _highlighted={{
                                    color: "neon.purple",
                                    bg: "glass.bgHover",
                                    filter: "drop-shadow(0 0 8px rgba(212, 0, 255, 0.45))",
                                }}
                                _selected={{
                                    color: "neon.blue",
                                    bg: "glass.bgHover",
                                    filter: "drop-shadow(0 0 8px rgba(0, 212, 255, 0.55))",
                                    fontWeight: "semibold",
                                }}
                                _focus={{
                                    outline: "none",
                                    color: "neon.blue",
                                    bg: "glass.bgHover",
                                    boxShadow: "0 0 0 1px rgba(0, 212, 255, 0.35)",
                                }}
                            >
                                {option.label}
                            </SelectItem>
                        ))}
                    </SelectContent>
                </SelectPositioner>
            </SelectRoot>
        </Field.Root>
    );
}
