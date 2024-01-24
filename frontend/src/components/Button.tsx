import {
  Box, Button, Icon, Text
} from "@chakra-ui/react";
import { LuSplit } from "react-icons/lu";
import { FcOk, FcHighPriority } from "react-icons/fc";

export function SplitButton({ isDisabled, onClick }) {
  return (
    <>
      <Box className="split-button-container">
        <Box
          className={`split-button ${isDisabled ? 'disabled' : ''}`}
          onClick={isDisabled ? null : onClick}
        >
          <Box>
            <Icon as={LuSplit} />
          </Box>
          <Text>
            Split
          </Text>
        </Box>
      </Box>
    </>
  );
}

export function ImportButton({ content, isDisabled, isLoading, onClick }) {
  return (
    <>
      <Button
        className="import-button-container"
        isDisabled={isDisabled}
        isLoading={isLoading}
        onClick={onClick}
      >
        {(() => {
          if (content === "Done") {
            return (
              <Icon as={FcOk} />
            )
          } else if (content === "Error") {
            return (
              <Icon as={FcHighPriority} />
            )
          }
        })()}
        <Text>{content}</Text>
      </Button>
    </>
  );
}
