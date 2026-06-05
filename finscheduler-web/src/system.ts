import {
    createSystem,
    defaultConfig,
    defineConfig,
    defineGlobalStyles,
    defineRecipe,
    defineSemanticTokens,
    defineSlotRecipe,
    defineTextStyles,
    defineTokens,
} from '@chakra-ui/react';

const darkOnly = (token: string) => token;

const globalCss = defineGlobalStyles({
    html: {
        bg: 'app.bg',
        color: 'fg',
        colorScheme: 'dark',
        lineHeight: '1.5',
        colorPalette: 'gray',
    },
    body: {
        margin: '0',
        minWidth: '320px',
        minHeight: '100vh',
        bg: 'app.bg',
        color: 'fg',
        fontFamily: 'body',
    },
    '#root': {
        minHeight: '100vh',
    },
    '*::placeholder, *[data-placeholder]': {
        color: 'fg.subtle',
    },
    '*::selection': {
        bg: 'app.accentSoft',
        color: 'fg',
    },
});

const textStyles = defineTextStyles({
    '2xs': {value: {fontSize: '0.6875rem', lineHeight: '1rem'}},
    xs: {value: {fontSize: '0.75rem', lineHeight: '1.125rem'}},
    sm: {value: {fontSize: '0.875rem', lineHeight: '1.375rem'}},
    md: {value: {fontSize: '1rem', lineHeight: '1.5rem'}},
    lg: {value: {fontSize: '1.125rem', lineHeight: '1.75rem'}},
    xl: {value: {fontSize: '1.25rem', lineHeight: '1.875rem'}},
    '2xl': {value: {fontSize: '1.5rem', lineHeight: '2.125rem', letterSpacing: '-0.02em'}},
    '3xl': {value: {fontSize: '1.875rem', lineHeight: '2.5rem', letterSpacing: '-0.025em'}},
    label: {
        value: {
            fontSize: '0.875rem',
            lineHeight: '1.25rem',
            fontWeight: '600',
            letterSpacing: '0.01em',
        },
    },
    eyebrow: {
        value: {
            fontSize: '0.75rem',
            lineHeight: '1rem',
            fontWeight: '600',
            letterSpacing: '0.08em',
            textTransform: 'uppercase',
        },
    },
    metric: {
        value: {
            fontSize: '1.875rem',
            lineHeight: '2.25rem',
            fontWeight: '700',
            letterSpacing: '-0.03em',
        },
    },
    sectionTitle: {
        value: {
            fontSize: '1.125rem',
            lineHeight: '1.625rem',
            fontWeight: '600',
            letterSpacing: '-0.015em',
        },
    },
});

const buttonRecipe = defineRecipe({
    className: 'chakra-button',
    base: {
        display: 'inline-flex',
        appearance: 'none',
        alignItems: 'center',
        justifyContent: 'center',
        userSelect: 'none',
        position: 'relative',
        borderRadius: 'l2',
        whiteSpace: 'nowrap',
        verticalAlign: 'middle',
        borderWidth: '1px',
        borderColor: 'transparent',
        cursor: 'button',
        flexShrink: '0',
        outline: '0',
        lineHeight: '1.2',
        isolation: 'isolate',
        fontWeight: '600',
        color: 'fg',
        transitionProperty: 'common',
        transitionDuration: 'moderate',
        backdropFilter: 'blur(18px)',
        _icon: {
            flexShrink: '0',
        },
        _disabled: {
            opacity: '0.45',
            cursor: 'not-allowed',
            boxShadow: 'none',
        },
    },
    variants: {
        size: {
            xs: {
                h: '8',
                minW: '8',
                px: '3',
                gap: '1.5',
                textStyle: 'xs',
                _icon: {width: '4', height: '4'},
            },
            sm: {
                h: '9',
                minW: '9',
                px: '3.5',
                gap: '2',
                textStyle: 'sm',
                _icon: {width: '4', height: '4'},
            },
            md: {
                h: '10',
                minW: '10',
                px: '4',
                gap: '2',
                textStyle: 'sm',
                _icon: {width: '4.5', height: '4.5'},
            },
            lg: {
                h: '11',
                minW: '11',
                px: '5',
                gap: '2.5',
                textStyle: 'md',
                _icon: {width: '5', height: '5'},
            },
        },
        variant: {
            solid: {
                bg: 'app.accent',
                color: 'bg',
                borderColor: 'transparent',
                boxShadow: 'app.glowCyan',
                _hover: {
                    bg: 'app.info',
                    transform: 'translateY(-1px)',
                },
                _expanded: {
                    bg: 'app.info',
                },
            },
            surface: {
                bg: 'app.cardBg',
                color: 'fg',
                borderColor: 'app.cardBorder',
                boxShadow: 'card',
                _hover: {
                    bg: 'app.cardBgHover',
                    borderColor: 'app.cardBorderActive',
                },
                _expanded: {
                    bg: 'app.cardBgHover',
                    borderColor: 'app.cardBorderActive',
                },
            },
            outline: {
                bg: 'transparent',
                color: 'fg',
                borderColor: 'app.cardBorder',
                _hover: {
                    bg: 'app.cardBg',
                    borderColor: 'app.cardBorderActive',
                },
                _expanded: {
                    bg: 'app.cardBg',
                    borderColor: 'app.cardBorderActive',
                },
            },
            ghost: {
                bg: 'transparent',
                color: 'fg.muted',
                _hover: {
                    bg: 'app.accentSoft',
                    color: 'fg',
                },
                _expanded: {
                    bg: 'app.accentSoft',
                    color: 'fg',
                },
            },
            selected: {
                bg: 'app.accentSoft',
                color: 'app.accent',
                borderColor: 'app.cardBorderActive',
                boxShadow: 'app.glowCyan',
                _hover: {
                    bg: 'app.accentSoft',
                },
                _expanded: {
                    bg: 'app.accentSoft',
                },
            },
        },
    },
    defaultVariants: {
        size: 'md',
        variant: 'solid',
    },
});

const inputRecipe = defineRecipe({
    className: 'chakra-input',
    base: {
        width: '100%',
        minWidth: '0',
        outline: '0',
        position: 'relative',
        appearance: 'none',
        textAlign: 'start',
        borderRadius: 'l2',
        height: 'var(--input-height)',
        minW: 'var(--input-height)',
        color: 'fg',
        _placeholder: {
            color: 'fg.subtle',
        },
        _disabled: {
            opacity: '0.45',
            cursor: 'not-allowed',
        },
        _invalid: {
            borderColor: 'border.error',
            focusRingColor: 'border.error',
        },
    },
    variants: {
        size: {
            xs: {
                textStyle: 'xs',
                px: '3',
                '--input-height': 'sizes.8',
            },
            sm: {
                textStyle: 'sm',
                px: '3',
                '--input-height': 'sizes.9',
            },
            md: {
                textStyle: 'sm',
                px: '3.5',
                '--input-height': 'sizes.10',
            },
            lg: {
                textStyle: 'md',
                px: '4',
                '--input-height': 'sizes.11',
            },
        },
        variant: {
            outline: {
                bg: 'app.cardBg',
                borderWidth: '1px',
                borderColor: 'app.cardBorder',
                boxShadow: 'inset',
                focusVisibleRing: 'inside',
                focusRingColor: 'app.accent',
                _hover: {
                    borderColor: 'app.cardBorderActive',
                },
                _focusVisible: {
                    bg: 'app.cardBgHover',
                    borderColor: 'app.cardBorderActive',
                    boxShadow: 'app.glowCyan',
                },
            },
            subtle: {
                bg: 'app.sidebarBg',
                borderWidth: '1px',
                borderColor: 'transparent',
                focusVisibleRing: 'inside',
                focusRingColor: 'app.accent',
                _hover: {
                    borderColor: 'app.cardBorder',
                },
                _focusVisible: {
                    bg: 'app.cardBg',
                    borderColor: 'app.cardBorderActive',
                    boxShadow: 'app.glowCyan',
                },
            },
            flushed: {
                bg: 'transparent',
                borderBottomWidth: '1px',
                borderBottomColor: 'app.cardBorder',
                borderRadius: '0',
                px: '0',
                _focusVisible: {
                    borderColor: 'app.accent',
                    boxShadow: '0px 1px 0px 0px var(--chakra-colors-app-accent)',
                },
            },
        },
    },
    defaultVariants: {
        size: 'md',
        variant: 'outline',
    },
});

const badgeRecipe = defineRecipe({
    className: 'chakra-badge',
    base: {
        display: 'inline-flex',
        alignItems: 'center',
        borderRadius: 'full',
        gap: '1',
        fontWeight: '600',
        fontVariantNumeric: 'tabular-nums',
        whiteSpace: 'nowrap',
        userSelect: 'none',
        backdropFilter: 'blur(14px)',
    },
    variants: {
        variant: {
            solid: {
                bg: 'app.accent',
                color: 'bg',
            },
            subtle: {
                bg: 'app.accentSoft',
                color: 'app.accent',
            },
            outline: {
                bg: 'transparent',
                color: 'fg.muted',
                borderWidth: '1px',
                borderColor: 'app.cardBorder',
            },
            surface: {
                bg: 'app.cardBg',
                color: 'fg',
                borderWidth: '1px',
                borderColor: 'app.cardBorder',
            },
        },
        size: {
            xs: {
                textStyle: '2xs',
                px: '2',
                minH: '5',
            },
            sm: {
                textStyle: 'xs',
                px: '2.5',
                minH: '6',
            },
            md: {
                textStyle: 'sm',
                px: '3',
                minH: '7',
            },
        },
    },
    defaultVariants: {
        variant: 'surface',
        size: 'sm',
    },
});

const cardRecipe = defineSlotRecipe({
    className: 'chakra-card',
    slots: ['root', 'header', 'body', 'footer', 'title', 'description'],
    base: {
        root: {
            display: 'flex',
            flexDirection: 'column',
            position: 'relative',
            minWidth: '0',
            wordWrap: 'break-word',
            borderRadius: 'l3',
            color: 'fg',
            textAlign: 'start',
            bg: 'app.cardBg',
            borderWidth: '1px',
            borderColor: 'app.cardBorder',
            boxShadow: 'card',
            backdropFilter: 'blur(20px)',
            transitionProperty: 'common',
            transitionDuration: 'moderate',
        },
        title: {
            fontWeight: '600',
            textStyle: 'sectionTitle',
            color: 'fg',
        },
        description: {
            color: 'fg.muted',
            textStyle: 'sm',
        },
        header: {
            paddingInline: 'var(--card-padding)',
            paddingTop: 'var(--card-padding)',
            display: 'flex',
            flexDirection: 'column',
            gap: '2',
        },
        body: {
            padding: 'var(--card-padding)',
            flex: '1',
            display: 'flex',
            flexDirection: 'column',
        },
        footer: {
            display: 'flex',
            alignItems: 'center',
            gap: '3',
            paddingInline: 'var(--card-padding)',
            paddingBottom: 'var(--card-padding)',
        },
    },
    variants: {
        size: {
            sm: {
                root: {
                    '--card-padding': 'spacing.5',
                },
                title: {
                    textStyle: 'lg',
                },
            },
            md: {
                root: {
                    '--card-padding': 'spacing.6',
                },
                title: {
                    textStyle: 'xl',
                },
            },
            lg: {
                root: {
                    '--card-padding': 'spacing.7',
                },
                title: {
                    textStyle: '2xl',
                },
            },
        },
        variant: {
            outline: {
                root: {
                    bg: 'app.cardBg',
                    borderColor: 'app.cardBorder',
                },
            },
            elevated: {
                root: {
                    bg: 'app.cardBgHover',
                    borderColor: 'app.cardBorder',
                    boxShadow: 'lg',
                },
            },
            selected: {
                root: {
                    bg: 'app.cardBgHover',
                    borderColor: 'app.cardBorderActive',
                    boxShadow: 'app.glowViolet',
                },
            },
        },
    },
    defaultVariants: {
        variant: 'outline',
        size: 'md',
    },
});

const tagRecipe = defineSlotRecipe({
    className: 'chakra-tag',
    slots: ['root', 'label', 'closeTrigger', 'startElement', 'endElement'],
    base: {
        root: {
            display: 'inline-flex',
            alignItems: 'center',
            maxWidth: '100%',
            userSelect: 'none',
            borderRadius: 'full',
            borderWidth: '1px',
            borderColor: 'app.cardBorder',
            bg: 'app.cardBg',
            color: 'fg',
            backdropFilter: 'blur(14px)',
            boxShadow: 'inset',
        },
        label: {
            lineClamp: '1',
            fontWeight: '600',
        },
        closeTrigger: {
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            outline: '0',
            borderRadius: 'full',
            color: 'currentColor',
        },
        startElement: {
            flexShrink: 0,
            boxSize: 'var(--tag-element-size)',
            ms: 'var(--tag-element-offset)',
            _icon: {boxSize: '100%'},
        },
        endElement: {
            flexShrink: 0,
            boxSize: 'var(--tag-element-size)',
            me: 'var(--tag-element-offset)',
            _icon: {boxSize: '100%'},
        },
    },
    variants: {
        size: {
            sm: {
                root: {
                    px: '2.5',
                    minH: '6',
                    gap: '1.5',
                    '--tag-element-size': 'spacing.3.5',
                    '--tag-element-offset': '-2px',
                },
                label: {
                    textStyle: 'xs',
                },
            },
            md: {
                root: {
                    px: '3',
                    minH: '7',
                    gap: '1.5',
                    '--tag-element-size': 'spacing.4',
                    '--tag-element-offset': '-3px',
                },
                label: {
                    textStyle: 'sm',
                },
            },
        },
        variant: {
            subtle: {
                root: {
                    bg: 'app.accentSoft',
                    color: 'app.accent',
                    borderColor: 'transparent',
                },
            },
            solid: {
                root: {
                    bg: 'app.accent',
                    color: 'bg',
                    borderColor: 'transparent',
                    boxShadow: 'app.glowCyan',
                },
            },
            outline: {
                root: {
                    bg: 'transparent',
                    color: 'fg.muted',
                    borderColor: 'app.cardBorder',
                },
            },
            surface: {
                root: {
                    bg: 'app.cardBg',
                    color: 'fg',
                    borderColor: 'app.cardBorder',
                },
            },
        },
    },
    defaultVariants: {
        size: 'md',
        variant: 'surface',
    },
});

const selectRecipe = defineSlotRecipe({
    className: 'chakra-select',
    slots: [
        'root',
        'trigger',
        'indicatorGroup',
        'indicator',
        'content',
        'item',
        'itemText',
        'itemGroup',
        'itemGroupLabel',
        'label',
        'valueText',
        'clearTrigger',
        'control',
    ],
    base: {
        root: {
            display: 'flex',
            flexDirection: 'column',
            gap: '1.5',
            width: 'full',
        },
        trigger: {
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            width: 'full',
            minH: 'var(--select-trigger-height)',
            '--input-height': 'var(--select-trigger-height)',
            px: 'var(--select-trigger-padding-x)',
            borderRadius: 'l2',
            userSelect: 'none',
            textAlign: 'start',
            borderWidth: '1px',
            borderColor: 'app.cardBorder',
            bg: 'app.cardBg',
            color: 'fg',
            boxShadow: 'inset',
            backdropFilter: 'blur(18px)',
            focusVisibleRing: 'inside',
            focusRingColor: 'app.accent',
            _placeholderShown: {
                color: 'fg.subtle',
            },
            _disabled: {
                opacity: '0.45',
            },
            _invalid: {
                borderColor: 'border.error',
            },
        },
        indicatorGroup: {
            display: 'flex',
            alignItems: 'center',
            gap: '1',
            pos: 'absolute',
            insetEnd: '0',
            top: '0',
            bottom: '0',
            px: 'var(--select-trigger-padding-x)',
            pointerEvents: 'none',
        },
        indicator: {
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: 'fg.muted',
        },
        content: {
            background: 'app.cardBgHover',
            display: 'flex',
            flexDirection: 'column',
            borderRadius: 'l2',
            outline: '0',
            maxH: '96',
            overflowY: 'auto',
            boxShadow: 'lg',
            borderWidth: '1px',
            borderColor: 'app.cardBorder',
            backdropFilter: 'blur(20px)',
        },
        item: {
            position: 'relative',
            userSelect: 'none',
            display: 'flex',
            alignItems: 'center',
            gap: '2',
            cursor: 'option',
            justifyContent: 'space-between',
            flex: '1',
            textAlign: 'start',
            borderRadius: 'l1',
            color: 'fg.muted',
            _highlighted: {
                bg: 'app.accentSoft',
                color: 'fg',
            },
            _selected: {
                bg: 'app.accentSoft',
                color: 'app.accent',
                boxShadow: 'app.glowCyan',
            },
            _disabled: {
                pointerEvents: 'none',
                opacity: '0.45',
            },
            _icon: {
                width: '4',
                height: '4',
            },
        },
        control: {
            pos: 'relative',
        },
        itemText: {
            flex: '1',
        },
        itemGroup: {
            _first: {mt: '0'},
        },
        itemGroupLabel: {
            py: '1',
            fontWeight: '600',
            color: 'fg.subtle',
        },
        label: {
            fontWeight: '600',
            userSelect: 'none',
            textStyle: 'label',
            color: 'fg.muted',
        },
        valueText: {
            lineClamp: '1',
            maxW: '80%',
        },
        clearTrigger: {
            color: 'fg.muted',
            pointerEvents: 'auto',
            borderRadius: 'full',
        },
    },
    variants: {
        variant: {
            outline: {
                trigger: {
                    bg: 'app.cardBg',
                    borderColor: 'app.cardBorder',
                    _expanded: {
                        bg: 'app.cardBgHover',
                        borderColor: 'app.cardBorderActive',
                        boxShadow: 'app.glowCyan',
                    },
                },
            },
            subtle: {
                trigger: {
                    bg: 'app.sidebarBg',
                    borderColor: 'transparent',
                    _expanded: {
                        bg: 'app.cardBg',
                        borderColor: 'app.cardBorder',
                    },
                },
            },
            ghost: {
                trigger: {
                    bg: 'transparent',
                    borderColor: 'transparent',
                    _expanded: {
                        bg: 'app.cardBg',
                        borderColor: 'app.cardBorder',
                    },
                },
            },
        },
        size: {
            sm: {
                root: {
                    '--select-trigger-height': 'sizes.9',
                    '--select-trigger-padding-x': 'spacing.3',
                },
                content: {
                    p: '1',
                    textStyle: 'sm',
                },
                trigger: {
                    textStyle: 'sm',
                    gap: '1.5',
                },
                item: {
                    py: '1.5',
                    px: '2',
                },
                itemGroupLabel: {
                    py: '1',
                    px: '2',
                },
                indicator: {
                    _icon: {
                        width: '4',
                        height: '4',
                    },
                },
            },
            md: {
                root: {
                    '--select-trigger-height': 'sizes.10',
                    '--select-trigger-padding-x': 'spacing.3.5',
                },
                content: {
                    p: '1.5',
                    textStyle: 'sm',
                },
                itemGroup: {
                    mt: '1.5',
                },
                item: {
                    py: '2',
                    px: '2.5',
                },
                itemGroupLabel: {
                    py: '1.5',
                    px: '2.5',
                },
                trigger: {
                    textStyle: 'sm',
                    gap: '2',
                },
                indicator: {
                    _icon: {
                        width: '4',
                        height: '4',
                    },
                },
            },
            lg: {
                root: {
                    '--select-trigger-height': 'sizes.11',
                    '--select-trigger-padding-x': 'spacing.4',
                },
                content: {
                    p: '1.5',
                    textStyle: 'md',
                },
                itemGroup: {
                    mt: '2',
                },
                item: {
                    py: '2',
                    px: '3',
                },
                itemGroupLabel: {
                    py: '2',
                    px: '3',
                },
                trigger: {
                    textStyle: 'md',
                    gap: '2',
                },
                indicator: {
                    _icon: {
                        width: '5',
                        height: '5',
                    },
                },
            },
        },
    },
    defaultVariants: {
        size: 'md',
        variant: 'outline',
    },
});

const nativeSelectRecipe = defineSlotRecipe({
    className: 'chakra-native-select',
    slots: ['root', 'field', 'indicator'],
    base: {
        root: {
            height: 'fit-content',
            display: 'flex',
            width: '100%',
            position: 'relative',
        },
        field: {
            width: '100%',
            minWidth: '0',
            outline: '0',
            appearance: 'none',
            borderRadius: 'l2',
            height: 'var(--select-field-height)',
            color: 'fg',
            backdropFilter: 'blur(18px)',
            '& > option, & > optgroup': {
                bg: 'app.sidebarBg',
            },
            _placeholder: {
                color: 'fg.subtle',
            },
            _disabled: {
                opacity: '0.45',
            },
            _invalid: {
                borderColor: 'border.error',
                focusRingColor: 'border.error',
            },
            focusVisibleRing: 'inside',
            focusRingColor: 'app.accent',
            lineHeight: 'normal',
        },
        indicator: {
            position: 'absolute',
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            pointerEvents: 'none',
            top: '50%',
            transform: 'translateY(-50%)',
            height: '100%',
            color: 'fg.muted',
            _disabled: {
                opacity: '0.45',
            },
            _icon: {
                width: '1em',
                height: '1em',
            },
        },
    },
    variants: {
        variant: {
            outline: {
                field: {
                    bg: 'app.cardBg',
                    borderWidth: '1px',
                    borderColor: 'app.cardBorder',
                    boxShadow: 'inset',
                    _focusVisible: {
                        bg: 'app.cardBgHover',
                        borderColor: 'app.cardBorderActive',
                        boxShadow: 'app.glowCyan',
                    },
                },
            },
            subtle: {
                field: {
                    bg: 'app.sidebarBg',
                    borderWidth: '1px',
                    borderColor: 'transparent',
                    _focusVisible: {
                        bg: 'app.cardBg',
                        borderColor: 'app.cardBorderActive',
                        boxShadow: 'app.glowCyan',
                    },
                },
            },
            ghost: {
                field: {
                    bg: 'transparent',
                    borderWidth: '1px',
                    borderColor: 'transparent',
                    _focusVisible: {
                        bg: 'app.cardBg',
                        borderColor: 'app.cardBorder',
                    },
                },
            },
        },
        size: {
            sm: {
                root: {
                    '--select-field-height': 'sizes.9',
                },
                field: {
                    textStyle: 'sm',
                    ps: '3',
                    pe: '9',
                },
                indicator: {
                    textStyle: 'md',
                    insetEnd: '2.5',
                },
            },
            md: {
                root: {
                    '--select-field-height': 'sizes.10',
                },
                field: {
                    textStyle: 'sm',
                    ps: '3.5',
                    pe: '9',
                },
                indicator: {
                    textStyle: 'lg',
                    insetEnd: '3',
                },
            },
            lg: {
                root: {
                    '--select-field-height': 'sizes.11',
                },
                field: {
                    textStyle: 'md',
                    ps: '4',
                    pe: '10',
                },
                indicator: {
                    textStyle: 'xl',
                    insetEnd: '3.5',
                },
            },
        },
    },
    defaultVariants: {
        size: 'md',
        variant: 'outline',
    },
});

const config = defineConfig({
    globalCss,
    theme: {
        tokens: defineTokens({
            colors: {
                navy: {
                    950: {value: '#050816'},
                    900: {value: '#081020'},
                    850: {value: '#0C1527'},
                    800: {value: '#101C31'},
                    750: {value: '#16263D'},
                },
                cyan: {
                    300: {value: '#6BE9FF'},
                    400: {value: '#20D0FF'},
                    500: {value: '#00B8D9'},
                },
                violet: {
                    300: {value: '#B7ABFF'},
                    400: {value: '#8F78FF'},
                    500: {value: '#6C55F4'},
                },
                green: {
                    300: {value: '#7DEAB4'},
                    400: {value: '#35D78A'},
                    500: {value: '#1FAE68'},
                },
                amber: {
                    300: {value: '#FFD38A'},
                    400: {value: '#FFB030'},
                    500: {value: '#E08A12'},
                },
                red: {
                    300: {value: '#FFA59B'},
                    400: {value: '#FF7060'},
                    500: {value: '#E14F42'},
                },
                glass: {
                    bg: {value: 'rgba(12, 21, 39, 0.78)'},
                    bgHover: {value: 'rgba(16, 29, 49, 0.92)'},
                    border: {value: 'rgba(32, 208, 255, 0.18)'},
                    borderStrong: {value: 'rgba(143, 120, 255, 0.34)'},
                },
                bg: {
                    base: {value: '{colors.navy.950}'},
                    layer1: {value: '{colors.navy.900}'},
                    layer2: {value: '{colors.navy.850}'},
                    accent: {value: '{colors.navy.800}'},
                },
                ink: {
                    primary: {value: '#F4F7FF'},
                    secondary: {value: 'rgba(214, 224, 241, 0.9)'},
                    muted: {value: 'rgba(186, 199, 223, 0.84)'},
                    disabled: {value: 'rgba(116, 131, 160, 0.72)'},
                    placeholder: {value: 'rgba(150, 167, 196, 0.78)'},
                },
                neon: {
                    blue: {value: '{colors.cyan.400}'},
                    cyan: {value: '{colors.cyan.300}'},
                    pink: {value: '{colors.violet.300}'},
                    purple: {value: '{colors.violet.400}'},
                    green: {value: '{colors.green.400}'},
                    yellow: {value: '{colors.amber.400}'},
                    shadow: {value: 'rgba(32, 208, 255, 0.34)'},
                },
                btc: {value: '{colors.cyan.400}'},
                eth: {value: '{colors.violet.400}'},
                ltc: {value: '{colors.green.400}'},
            },
            fonts: {
                heading: {value: '"Segoe UI", "Inter", sans-serif'},
                body: {value: '"Segoe UI", "Inter", sans-serif'},
                mono: {value: '"JetBrains Mono", "SFMono-Regular", monospace'},
            },
            fontWeights: {
                normal: {value: '400'},
                medium: {value: '500'},
                semibold: {value: '600'},
                bold: {value: '700'},
            },
            lineHeights: {
                tight: {value: '1.15'},
                snug: {value: '1.3'},
                normal: {value: '1.5'},
                relaxed: {value: '1.65'},
            },
            letterSpacings: {
                tighter: {value: '-0.03em'},
                tight: {value: '-0.02em'},
                normal: {value: '0'},
                wide: {value: '0.04em'},
            },
            spacing: {
                gutter: {value: '1.5rem'},
                panel: {value: '1.25rem'},
                section: {value: '1.75rem'},
                dense: {value: '0.75rem'},
            },
            radii: {
                l1: {value: '12px'},
                l2: {value: '18px'},
                l3: {value: '24px'},
                sm: {value: '12px'},
                md: {value: '18px'},
                lg: {value: '24px'},
                xl: {value: '32px'},
                full: {value: '9999px'},
            },
            shadows: {
                xs: {value: '0 4px 12px rgba(3, 8, 20, 0.22)'},
                sm: {value: '0 10px 24px rgba(3, 8, 20, 0.28)'},
                md: {value: '0 16px 36px rgba(3, 8, 20, 0.34)'},
                lg: {value: '0 24px 52px rgba(3, 8, 20, 0.42)'},
                xl: {value: '0 30px 72px rgba(3, 8, 20, 0.5)'},
                '2xl': {value: '0 40px 96px rgba(3, 8, 20, 0.56)'},
                inner: {value: 'inset 0 1px 0 rgba(255, 255, 255, 0.03)'},
                inset: {value: 'inset 0 0 0 1px rgba(255, 255, 255, 0.04)'},
                panel: {value: '0 18px 48px rgba(3, 8, 20, 0.36)'},
                glowCyan: {
                    value: '0 0 0 1px rgba(32, 208, 255, 0.18), 0 0 24px rgba(32, 208, 255, 0.18)',
                },
                glowViolet: {
                    value: '0 0 0 1px rgba(143, 120, 255, 0.22), 0 0 28px rgba(143, 120, 255, 0.18)',
                },
            },
            blurs: {
                soft: {value: '16px'},
                panel: {value: '20px'},
            },
        }),
        semanticTokens: defineSemanticTokens({
            colors: {
                bg: {
                    DEFAULT: {value: darkOnly('{colors.navy.950}')},
                    subtle: {value: darkOnly('{colors.navy.900}')},
                    muted: {value: darkOnly('{colors.navy.850}')},
                    emphasized: {value: darkOnly('{colors.navy.800}')},
                    inverted: {value: darkOnly('{colors.ink.primary}')},
                    panel: {value: darkOnly('{colors.glass.bg}')},
                    error: {value: darkOnly('{colors.red.500}')},
                    warning: {value: darkOnly('{colors.amber.500}')},
                    success: {value: darkOnly('{colors.green.500}')},
                    info: {value: darkOnly('{colors.cyan.500}')},
                },
                fg: {
                    DEFAULT: {value: darkOnly('{colors.ink.primary}')},
                    muted: {value: darkOnly('{colors.ink.secondary}')},
                    subtle: {value: darkOnly('{colors.ink.muted}')},
                    inverted: {value: darkOnly('{colors.navy.950}')},
                    error: {value: darkOnly('{colors.red.300}')},
                    warning: {value: darkOnly('{colors.amber.300}')},
                    success: {value: darkOnly('{colors.green.300}')},
                    info: {value: darkOnly('{colors.cyan.300}')},
                },
                border: {
                    DEFAULT: {value: darkOnly('{colors.glass.border}')},
                    muted: {value: darkOnly('rgba(255, 255, 255, 0.04)')},
                    subtle: {value: darkOnly('rgba(255, 255, 255, 0.02)')},
                    emphasized: {value: darkOnly('{colors.glass.borderStrong}')},
                    inverted: {value: darkOnly('{colors.ink.primary}')},
                    error: {value: darkOnly('{colors.red.400}')},
                    warning: {value: darkOnly('{colors.amber.400}')},
                    success: {value: darkOnly('{colors.green.400}')},
                    info: {value: darkOnly('{colors.cyan.400}')},
                },
                app: {
                    bg: {value: darkOnly('{colors.navy.950}')},
                    sidebarBg: {value: darkOnly('rgba(8, 16, 32, 0.88)')},
                    topbarBg: {value: darkOnly('rgba(8, 16, 32, 0.72)')},
                    cardBg: {value: darkOnly('{colors.glass.bg}')},
                    cardBgHover: {value: darkOnly('{colors.glass.bgHover}')},
                    cardBorder: {value: darkOnly('{colors.glass.border}')},
                    cardBorderActive: {value: darkOnly('rgba(143, 120, 255, 0.42)')},
                    text: {value: darkOnly('{colors.ink.primary}')},
                    textMuted: {value: darkOnly('{colors.ink.secondary}')},
                    textSubtle: {value: darkOnly('{colors.ink.muted}')},
                    accent: {value: darkOnly('{colors.cyan.400}')},
                    accentSoft: {value: darkOnly('rgba(32, 208, 255, 0.16)')},
                    accentViolet: {value: darkOnly('{colors.violet.400}')},
                    positive: {value: darkOnly('{colors.green.400}')},
                    negative: {value: darkOnly('{colors.red.400}')},
                    warning: {value: darkOnly('{colors.amber.400}')},
                    info: {value: darkOnly('{colors.cyan.300}')},
                    glowCyan: {value: darkOnly('rgba(32, 208, 255, 0.34)')},
                    glowViolet: {value: darkOnly('rgba(143, 120, 255, 0.3)')},
                    sidebarItemHoverStart: {value: darkOnly('rgba(32, 208, 255, 0.12)')},
                    sidebarItemHoverEnd: {value: darkOnly('rgba(143, 120, 255, 0.1)')},
                    sidebarItemHoverBorder: {value: darkOnly('rgba(96, 165, 255, 0.28)')},
                },
                background: {value: darkOnly('{colors.navy.950}')},
                surface: {value: darkOnly('{colors.glass.bg}')},
                surfaceHover: {value: darkOnly('{colors.glass.bgHover}')},
                borderStrong: {value: darkOnly('{colors.glass.borderStrong}')},
                textPrimary: {value: darkOnly('{colors.ink.primary}')},
                textSecondary: {value: darkOnly('{colors.ink.secondary}')},
                textMuted: {value: darkOnly('{colors.ink.muted}')},
                textDisabled: {value: darkOnly('{colors.ink.disabled}')},
            },
            shadows: {
                card: {
                    value: darkOnly(
                        '0 18px 48px rgba(3, 8, 20, 0.36), inset 0 1px 0 rgba(255, 255, 255, 0.03)',
                    ),
                },
                neon: {
                    value: darkOnly(
                        '0 0 0 1px rgba(32, 208, 255, 0.18), 0 0 24px rgba(32, 208, 255, 0.18)',
                    ),
                },
                app: {
                    glowCyan: {
                        value: darkOnly(
                            '0 0 0 1px rgba(32, 208, 255, 0.18), 0 0 24px rgba(32, 208, 255, 0.18)',
                        ),
                    },
                    glowViolet: {
                        value: darkOnly(
                            '0 0 0 1px rgba(143, 120, 255, 0.22), 0 0 28px rgba(143, 120, 255, 0.18)',
                        ),
                    },
                    sidebarItemHover: {
                        value: darkOnly(
                            '0 0 0 1px rgba(96, 165, 255, 0.2), inset 0 1px 0 rgba(255, 255, 255, 0.03)',
                        ),
                    },
                },
            },
        }),
        textStyles,
        recipes: {
            badge: badgeRecipe,
            button: buttonRecipe,
            input: inputRecipe,
        },
        slotRecipes: {
            card: cardRecipe,
            nativeSelect: nativeSelectRecipe,
            select: selectRecipe,
            tag: tagRecipe,
        },
    },
});

export const system = createSystem(defaultConfig, config);
