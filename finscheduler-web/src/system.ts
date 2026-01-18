import {createSystem, defaultConfig} from "@chakra-ui/react"

export const system = createSystem(defaultConfig, {
    theme: {
        tokens: {
            fonts: {
                body: {value: " "},
            },
            fontWeights: {
                normal: {value: 400},
                medium: {value: 500},
                semibold: {value: 600},
                bold: {value: 700},
            },

            colors: {
                neon: {
                    blue: {value: "rgba(0,212,255)"},
                    cyan: {value: "#00fff0"},
                    pink: {value: "#ff00c8"},
                    purple: {value: "#c800ff"},
                    green: {value: "#00ffa3"},
                    yellow: {value: "#c8ff00"},
                    shadow: {value: "rgba(0, 212, 255, 0.9)"},
                },

                bg: {
                    base: {value: "#0b0e1a"},
                    layer1: {value: "#11152b"},
                    layer2: {value: "#141a38"},
                    accent: {value: "#1a1f40"},
                },

                glass: {
                    bg: {value: "rgba(17, 21, 43, 0.65)"},
                    bgHover: {value: "rgba(20, 26, 56, 0.8)"},
                    border: {value: "rgba(255, 255, 255, 0.1)"},
                    borderStrong: {value: "rgba(255, 255, 255, 0.18)"},
                },

                text: {
                    primary: {value: "#ffffff"},
                    secondary: {value: "rgba(255, 255, 255, 0.75)"},
                    muted: {value: "rgba(255, 255, 255, 0.55)"},
                    disabled: {value: "rgba(255, 255, 255, 0.35)"},
                },

                btc: {value: "#00d4ff"},
                eth: {value: "#ff00c8"},
                ltc: {value: "#00ffa3"},
            },

            radii: {
                sm: {value: "8px"},
                md: {value: "16px"},
                lg: {value: "24px"},
                xl: {value: "32px"},
                full: {value: "9999px"},
            },

            shadows: {
                sm: {value: "0 4px 20px rgba(0, 0, 0, 0.3)"},
                md: {value: "0 8px 32px rgba(0, 0, 0, 0.4)"},
                lg: {value: "0 16px 48px rgba(0, 0, 0, 0.5)"},
                glass: {value: "0 8px 32px rgba(0, 0, 0, 0.37)"},
                inner: {value: "inset 0 2px 10px rgba(0, 0, 0, 0.6)"},

                neonBlue: {value: "0 0 20px rgba(0, 212, 255, 0.4)"},
                neonPink: {value: "0 0 20px rgba(255, 0, 200, 0.4)"},
                neonPurple: {value: "0 0 20px rgba(200, 0, 255, 0.4)"},
                neonGreen: {value: "0 0 20px rgba(0, 255, 163, 0.4)"},
            },

            gradients: {
                cosmic: {
                    value: "radial-gradient(circle at 10% 20%, rgba(0, 212, 255, 0.15) 0%, transparent 20%), radial-gradient(circle at 90% 80%, rgba(255, 0, 200, 0.15) 0%, transparent 20%), radial-gradient(circle at 50% 50%, rgba(200, 0, 255, 0.08) 0%, transparent 30%)",
                },
                cardGlow: {
                    value: "linear-gradient(135deg, rgba(0, 212, 255, 0.15), rgba(255, 0, 200, 0.15))",
                },
            },

            blurs: {
                glass: {value: "blur(16px)"},
                glassLight: {value: "blur(10px)"},
            },
        },

        semanticTokens: {
            colors: {
                background: {value: {base: "{colors.bg.base}", _dark: "{colors.bg.base}"}},
                surface: {value: "{colors.glass.bg}"},
                surfaceHover: {value: "{colors.glass.bgHover}"},
                border: {value: "{colors.glass.border}"},
                borderStrong: {value: "{colors.glass.borderStrong}"},
                textPrimary: {value: "{colors.text.primary}"},
                textSecondary: {value: "{colors.text.secondary}"},
                textMuted: {value: "{colors.text.muted}"},
            },
            shadows: {
                card: {value: "{shadows.glass}"},
                neon: {value: "{shadows.neonBlue}"},
            },
        },
    },
})