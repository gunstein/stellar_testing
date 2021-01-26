import React, { useState, useContext } from "react";
import { RoutingContext, pagesMapping } from "../context/RoutingContext";
import selectedTileContext from "../context/SelectedTileContext";
import { makeStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import GalleryAppBar from "../components/GalleryAppBar";
import Button from "@material-ui/core/Button";
import {
  createOrder,
  getEventSourceSSEReceiveBigImageUrl,
  getBigFileUrlFromMemoAndDownloadKey,
} from "../services/order";
import {
  createAccount,
  getAccountBalance,
  sendPayment,
} from "../services/stellar";
import CircularProgress from "@material-ui/core/CircularProgress";
import HomeIcon from "@material-ui/icons/Home";
import Fab from "@material-ui/core/Fab";
import ImageDialog from "../components/ImageDialog";

const useStyles = makeStyles((theme) => ({
  appBar: {
    position: "relative",
  },
  layout: {
    width: "auto",
    marginLeft: theme.spacing(2),
    marginRight: theme.spacing(2),
    [theme.breakpoints.up(600 + theme.spacing(2) * 2)]: {
      width: 600,
      marginLeft: "auto",
      marginRight: "auto",
    },
  },
  paper: {
    marginTop: theme.spacing(3),
    marginBottom: theme.spacing(3),
    padding: theme.spacing(2),
    [theme.breakpoints.up(600 + theme.spacing(3) * 2)]: {
      marginTop: theme.spacing(6),
      marginBottom: theme.spacing(6),
      padding: theme.spacing(3),
    },
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
}));

const Buy = () => {
  const { setPage } = useContext(RoutingContext);
  const { selectedTile } = useContext(selectedTileContext);
  const [targetAccount, setTargetAccount] = useState("");
  const [targetMemo, setTargetMemo] = useState("");
  const [targetDownloadKey, setTargetDownloadKey] = useState("");
  const [sourceAccountPublicKey, setSourceAccountPublicKey] = useState("");
  const [sourceAccountSecret, setSourceAccountSecret] = useState("");
  const [sourceAccountBalance, setSourceAccountBalance] = useState(0);
  const [orderMode, setOrderMode] = useState(true);
  const [createAccountMode, setCreateAccountMode] = useState(false);
  const [payMode, setPayMode] = useState(false);
  const [spinnerCreateAccount, setSpinnerCreateAccount] = useState(false);
  const [spinnerPay, setSpinnerPay] = useState(false);
  const [sourceKeyPair, setSourceKeyPair] = useState(null);
  const [bigImageUrl, setBigImageUrl] = useState("");
  //const [receiveBigImageMode, setReceiveBigImageMode] = useState(false);

  const classes = useStyles();

  var evSource = null;

  const handleOrder = () => {
    createOrder({ artid: selectedTile.artid }).then((order_response) => {
      console.log("order_response: ", order_response);
      if (
        order_response.hasOwnProperty("data") &&
        order_response.data.hasOwnProperty("account") &&
        order_response.data.hasOwnProperty("memo")
      ) {
        setTargetAccount(order_response.data.account);
        setTargetMemo(order_response.data.memo);
        setTargetDownloadKey(order_response.data.download_key);
        setOrderMode(false);
        setCreateAccountMode(true);
      }
    });
  };

  const handleCreateAccount = () => {
    setSpinnerCreateAccount(true);
    createAccount().then((account_response) => {
      setSourceKeyPair(account_response);
      setSourceAccountPublicKey(account_response.publicKey());
      setSourceAccountSecret(account_response.secret());
      handleGetAccountBalance(account_response.publicKey());
      setCreateAccountMode(false);
      setPayMode(true);
      setSpinnerCreateAccount(false);
    });
  };

  const handleGetAccountBalance = (sourceAccountPublicKey) => {
    getAccountBalance(sourceAccountPublicKey).then(
      (account_balance_response) => {
        setSourceAccountBalance(account_balance_response.balance);
      }
    );
  };

  const handleReceiveSSEMessage = (event) => {
    console.log("handleReceiveSSEMessage: ", event);
    if (event.data === targetMemo) {
      console.log("before getBigFileUrlFromMemoAndDownloadKey");
      getBigFileUrlFromMemoAndDownloadKey(targetMemo, targetDownloadKey).then(
        (data) => {
          console.log("data: ", data);
          setBigImageUrl(data.data.big_file_url);
          evSource.removeEventListener("message", handleReceiveSSEMessage);
          evSource.close();
          handleGetAccountBalance(sourceAccountPublicKey);
          setPayMode(false);
          setSpinnerPay(false);
        }
      );
    }
  };

  const handleSSEError = (error) => {
    console.error("EventSource failed:", error);
  };

  const handlePay = () => {
    setSpinnerPay(true);

    evSource = getEventSourceSSEReceiveBigImageUrl();
    evSource.onerror = handleSSEError;
    evSource.addEventListener("message", handleReceiveSSEMessage);

    sendPayment(
      sourceKeyPair,
      targetAccount,
      targetMemo,
      selectedTile.price
    ).then();
  };

  const handleGallery = () => {
    setPage(pagesMapping.gallery);
  };

  const handleCloseImageDialog = () => {
    setBigImageUrl("");
  };

  return (
    <React.Fragment>
      {bigImageUrl === "" ? (
        <main className={classes.layout}>
          <GalleryAppBar></GalleryAppBar>
          <Paper className={classes.paper}>
            <Typography variant="h6" gutterBottom>
              Buy {selectedTile.title} by {selectedTile.artist}
            </Typography>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <TextField
                  id="target_account"
                  name="targetAccount"
                  label="Target Account"
                  fullWidth
                  InputProps={{
                    readOnly: true,
                  }}
                  value={targetAccount}
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  id="target_memo"
                  name="targetMemo"
                  label="Target Memo"
                  fullWidth
                  value={targetMemo}
                  InputProps={{
                    readOnly: true,
                  }}
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  id="target_amount"
                  name="targetAmount"
                  label="Target Amount"
                  fullWidth
                  value={selectedTile.price}
                  InputProps={{
                    readOnly: true,
                  }}
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  id="target_downloadkey"
                  name="targetDownloadKey"
                  label="Target DownloadKey"
                  fullWidth
                  value={targetDownloadKey}
                  InputProps={{
                    readOnly: true,
                  }}
                />
              </Grid>
              {orderMode ? (
                <Grid item xs={12}>
                  <div className={classes.buttons}>
                    <Button
                      variant="contained"
                      color="primary"
                      className={classes.button}
                      onClick={handleOrder}
                    >
                      Order
                    </Button>
                  </div>
                </Grid>
              ) : null}
              {!orderMode ? (
                <React.Fragment>
                  <Grid item xs={12}>
                    <TextField
                      id="source_account_public_key"
                      name="sourceAccountPublicKey"
                      label="Source Account Public Key"
                      fullWidth
                      value={sourceAccountPublicKey}
                      InputProps={{
                        readOnly: true,
                      }}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <TextField
                      id="source_account_secret"
                      name="sourceAccountSecret"
                      label="Source Account Secret"
                      fullWidth
                      value={sourceAccountSecret}
                      InputProps={{
                        readOnly: true,
                      }}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <TextField
                      id="source_account_balance"
                      name="sourceAccountBalance"
                      label="Source Account Balance"
                      fullWidth
                      value={sourceAccountBalance}
                      InputProps={{
                        readOnly: true,
                      }}
                    />
                  </Grid>
                </React.Fragment>
              ) : null}
              {createAccountMode ? (
                <Grid item xs={12}>
                  <div className={classes.buttons}>
                    {spinnerCreateAccount ? <CircularProgress /> : null}
                    <Button
                      variant="contained"
                      color="primary"
                      className={classes.button}
                      onClick={handleCreateAccount}
                    >
                      Create and fund account
                    </Button>
                  </div>
                </Grid>
              ) : null}
              {payMode ? (
                <Grid item xs={12}>
                  <div className={classes.buttons}>
                    {spinnerPay ? <CircularProgress /> : null}
                    <Button
                      variant="contained"
                      color="primary"
                      className={classes.button}
                      onClick={handlePay}
                    >
                      Pay
                    </Button>
                  </div>
                </Grid>
              ) : null}
            </Grid>
            <Fab
              color="primary"
              aria-label="Gallery"
              onClick={() => handleGallery()}
            >
              <HomeIcon />
            </Fab>
          </Paper>
        </main>
      ) : null}
      {bigImageUrl !== "" ? (
        <ImageDialog
          handleCloseToParent={handleCloseImageDialog}
          bigImageUrl={bigImageUrl}
        ></ImageDialog>
      ) : null}
    </React.Fragment>
  );
};

export default Buy;
