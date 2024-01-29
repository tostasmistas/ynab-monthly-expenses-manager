import React, { useState, useEffect } from "react";
import { render } from "react-dom";
import {
  ChakraProvider, Alert, AlertIcon, AlertDescription, Box, Flex, Spinner
} from "@chakra-ui/react";

import "./index.css";
import { theme } from "./themes/theme"
import { Header } from "./components/Header"
import { SharedMonthlyExpensesCard, IndividualMonthlyExpensesCard } from "./components/MonthlyExpensesCard"
import { SplitButton, ImportButton } from "./components/Button"

import { backend } from "../wailsjs/go/models";
import { GetSharedMonthlyExpenses, CreateMonthlyExpensesTransactions } from "../wailsjs/go/backend/Backend";
import { EventsEmit, EventsOn } from "../wailsjs/runtime";

const App = () => {
  const [backendLoaded, setBackendLoaded] = useState(null)

  const [sharedMonthlyExpenses, setSharedMonthlyExpenses] = useState<backend.MonthlyExpenses>()
  const [individualMonthlyExpenses, setIndividualMonthlyExpenses] = useState<backend.MonthlyExpenses>()

  const [splitButtonDisabled, setSplitButtonDisabled] = useState(true)

  const [importButtonContent, setImportButtonContent] = useState("Import")
  const [importButtonDisabled, setImportButtonDisabled] = useState(true)
  const [importButtonLoading, setImportButtonLoading] = useState(false)

  useEffect(() => {
    EventsOn("backendSetupComplete", function(args?: any) {
      setTimeout(() => {
        setBackendLoaded(args);
      }, 1000);
    })
  }, []);

  useEffect(() => {
    GetSharedMonthlyExpenses().then(monthlyExpenses=> {
      setSharedMonthlyExpenses(monthlyExpenses);
    });
  }, []);

  useEffect(() => {
    EventsOn("sharedMonthlyExpensesSplit", function(args?: any) {
      setIndividualMonthlyExpenses(args);
      setImportButtonDisabled(false);
      if (importButtonContent !== "Import") {
        setImportButtonContent("Import");
      }
    })
  });

  const handleChange = (event) => {
    const { name, value } = event.target;

    setSharedMonthlyExpenses((previousSharedMonthlyExpenses) => ({
      ...previousSharedMonthlyExpenses,
      expenses: {
        ...previousSharedMonthlyExpenses.expenses,
        [name]: {
          ...previousSharedMonthlyExpenses.expenses[name],
          amount: value,
        },
      }
    }));

    setSplitButtonDisabled(false);
  };

  const splitSharedMonthlyExpenses = () => {
    EventsEmit("sharedMonthlyExpensesInput", sharedMonthlyExpenses);
  };

  const createMonthlyExpensesTransactions = () => {
    setSplitButtonDisabled(true);
    setImportButtonDisabled(true);
    setImportButtonLoading(true);

    CreateMonthlyExpensesTransactions(
      new backend.CombinedMonthlyExpenses({
        shared_monthly_expenses: sharedMonthlyExpenses,
        individual_monthly_expenses: individualMonthlyExpenses
      })
    ).then(response => {
      setTimeout(() => {
        setImportButtonLoading(false);
        if (response === true) {
          setImportButtonContent("Done");
        } else {
          setImportButtonContent("Error");
          setSplitButtonDisabled(false);
        }
      }, 1000);
    });
  }

  return (
    <>
      <ChakraProvider theme={theme}>
        <Box className="main-container">
          <Header/>
          <Flex className="body-container">
            <SharedMonthlyExpensesCard
              monthlyExpenses={sharedMonthlyExpenses}
              onChange={handleChange}
            />
            <Box className="buttons-container">
              <SplitButton
                isDisabled={splitButtonDisabled}
                onClick={splitSharedMonthlyExpenses}
              />
              <ImportButton
                content={importButtonContent}
                isDisabled={importButtonDisabled}
                isLoading={importButtonLoading}
                onClick={createMonthlyExpensesTransactions}
              />
            </Box>
            <IndividualMonthlyExpensesCard
              monthlyExpenses={individualMonthlyExpenses}
            />
          </Flex>
          {(() => {
            if (backendLoaded === null) {
              return (
                <Box className="overlay-container">
                  <Spinner
                    speed="0.65s"
                  />
                </Box>
              )
            } else if (backendLoaded === false) {
              return (
                <Box className="overlay-container">
                  <Alert status="error">
                    <AlertIcon />
                    <AlertDescription maxWidth='sm'>
                      Error setting up the application
                    </AlertDescription>
                  </Alert>
                </Box>
              )
            }
          })()}
        </Box>
      </ChakraProvider>
    </>
  );
};

render(<App/>, document.getElementById("root"));
