import {
    createListCollection,
    Field,
    SelectContent,
    SelectItem,
    SelectPositioner,
    SelectRoot,
    SelectTrigger,
    SelectValueText,
    Spinner,
} from "@chakra-ui/react";
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
                <SelectTrigger bg="bg.layer2" borderColor="glass.border">
                    <SelectValueText
                        placeholder={placeholder}
                        color="neon.blue"
                        overflow="hidden"
                        textOverflow="ellipsis"
                        whiteSpace="nowrap"
                    />
                </SelectTrigger>
                <SelectPositioner>
                    <SelectContent
                        bg="bg.layer1"
                        borderColor="glass.border"
                        zIndex="popover"
                        maxH={SELECT_CONTENT_MAX_HEIGHT}
                        overflowY="auto"
                        overscrollBehavior="contain"
                    >
                        {collection.items.map((option) => (
                            <SelectItem
                                item={option}
                                key={option.value}
                                _hover={{bg: "bg.layer3"}}
                                _selected={props.multiple ? {color: "neon.blue", fontWeight: "bold"} : undefined}
                                color="white"
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
