import {
  Avatar,
  Box,
  Card,
  CardBody,
  CardHeader,
  Flex,
  FormLabel,
  Icon,
  Input,
  InputGroup,
  InputLeftAddon,
  NumberInput,
  NumberInputField,
  Stack,
  StackDivider,
  Text
} from "@chakra-ui/react";
import { FcBusinesswoman, FcDepartment, FcHome, FcIdea, FcSimCard } from "react-icons/fc";
import { IoWater } from "react-icons/io5";

function MonthlyExpenseIcon({ categoryName }) {
  switch(categoryName) {
    case "Condominium":
      return <Icon
        as={FcDepartment}
      />;
    case "Electricity":
      return <Icon
        as={FcIdea}
        className="idea-icon"
      />;
    case "Water":
      return <Icon
        as={IoWater}
        color="blue.400"
      />;
    case "TV / Internet / Phone":
      return <Icon
        as={FcSimCard}
      />;
  }
}

function MonthlyExpenseInputLabel({ categoryName }) {
  return (
    <>
      <Flex className="expense-input-label-container">
        { MonthlyExpenseIcon({categoryName}) }
        <FormLabel>
          <Text>{categoryName}</Text>
        </FormLabel>
      </Flex>
    </>
  )
}

function MonthlyExpenseInput({ categoryName, amount, onChange }) {
  return (
    <>
      <Box className="expense-input-container">
        { MonthlyExpenseInputLabel({categoryName}) }
        <InputGroup size="md">
          <InputLeftAddon children="€"/>
          <NumberInput min={0} defaultValue={0} precision={2}>
            <NumberInputField
              placeholder="Enter an amount"
              name={categoryName}
              value={amount}
              onChange={onChange}
            />
          </NumberInput>
        </InputGroup>
      </Box>
    </>
  );
}

function MonthlyExpenseDisabledInput({ categoryName, amount }) {
  return (
    <>
      <Box className="expense-input-container">
        { MonthlyExpenseInputLabel({categoryName}) }
        <InputGroup size="md">
          <InputLeftAddon children="€"/>
          <Input
            isDisabled={true}
            placeholder="Enter an amount"
            name={categoryName}
            value={Math.max(parseFloat(amount) || 0, 0).toFixed(2)}
            onChange={null}
          />
        </InputGroup>
      </Box>
    </>
  );
}

export function SharedMonthlyExpensesCard({ monthlyExpenses, onChange }) {
  return (
    <>
      <Box className="expenses-card">
        <Card shadow="md">
          <CardHeader>
            <Flex>
              <Avatar icon={<FcHome />} className="shared-avatar" />
              <Text>Total Monthly Expenses</Text>
            </Flex>
          </CardHeader>
          <CardBody>
            <Stack divider={<StackDivider />} spacing="5">
              <Box>
                <MonthlyExpenseInput
                  categoryName="Condominium"
                  amount={monthlyExpenses?.expenses?.["Condominium"]?.amount || ""}
                  onChange={onChange}
                />
              </Box>
              <Box>
                <MonthlyExpenseInput
                  categoryName="Electricity"
                  amount={monthlyExpenses?.expenses?.["Electricity"]?.amount || ""}
                  onChange={onChange}
                />
              </Box>
              <Box>
                <MonthlyExpenseInput
                  categoryName="Water"
                  amount={monthlyExpenses?.expenses?.["Water"]?.amount || ""}
                  onChange={onChange}
                />
              </Box>
              <Box>
                <MonthlyExpenseInput
                  categoryName="TV / Internet / Phone"
                  amount={monthlyExpenses?.expenses?.["TV / Internet / Phone"]?.amount || ""}
                  onChange={onChange}
                />
              </Box>
            </Stack>
          </CardBody>
        </Card>
      </Box>
    </>
  );
}

export function IndividualMonthlyExpensesCard({ monthlyExpenses }) {
  return (
    <>
      <Box className="expenses-card">
        <Card shadow="md">
          <CardHeader>
            <Flex>
              <Avatar icon={<FcBusinesswoman />} className="individual-avatar" />
              <Text>Individual Share</Text>
            </Flex>
          </CardHeader>
          <CardBody>
            <Stack divider={<StackDivider />} spacing="5">
              <Box>
                <MonthlyExpenseDisabledInput
                  categoryName="Condominium"
                  amount={monthlyExpenses?.expenses?.["Condominium"]?.amount || ""}
                />
              </Box>
              <Box>
                <MonthlyExpenseDisabledInput
                  categoryName="Electricity"
                  amount={monthlyExpenses?.expenses?.["Electricity"]?.amount || ""}
                />
              </Box>
              <Box>
                <MonthlyExpenseDisabledInput
                  categoryName="Water"
                  amount={monthlyExpenses?.expenses?.["Water"]?.amount || ""}
                />
              </Box>
              <Box>
                <MonthlyExpenseDisabledInput
                  categoryName="TV / Internet / Phone"
                  amount={monthlyExpenses?.expenses?.["TV / Internet / Phone"]?.amount || ""}
                />
              </Box>
            </Stack>
          </CardBody>
        </Card>
      </Box>
    </>
  );
}
