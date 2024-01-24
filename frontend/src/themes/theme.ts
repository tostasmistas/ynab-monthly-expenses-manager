import {
  extendTheme
} from "@chakra-ui/react";

export const theme = extendTheme({
  styles: {
    global: () => ({
      html: {
        fontSize: "18px",
      },
      body: {
        bg: "#FAFAFA",
      },
    }),
  },
  fonts: {
    body: `'Open Sans', sans-serif`,
  },
  components: {
    NumberInput: {
      variants: {
        outline: {
          field: {
            _focusVisible: {
              boxShadow: "none !important",
            },
          }
        }
      }
    },
    Input: {
      baseStyle: {
        field: {
          _disabled: {
            opacity: 1,
          },
        },
      }
    },
    Button: {
      baseStyle: {
        _disabled: {
          opacity: 0.5,
        },
      }
    }
  }
})
