import {
  AbsoluteCenter, Box, Divider, Image
} from "@chakra-ui/react";
import YNABLogo from "../assets/images/ynab_logo.svg";

export function Header() {
  return (
    <>
      <Box className="header-container">
        <Divider />
        <AbsoluteCenter className="image-container">
          <Image src={YNABLogo} />
        </AbsoluteCenter>
      </Box>
    </>
  );
}
