import React, {Component} from 'react';
import './App.css';
import {createMuiTheme, MuiThemeProvider, withStyles} from "@material-ui/core/styles";
import {blueGrey} from '@material-ui/core/colors'
import AppBar from "@material-ui/core/AppBar";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import Toolbar from "@material-ui/core/Toolbar";
import {Menu as MenuIcon} from '@material-ui/icons';
import CssBaseline from "@material-ui/core/CssBaseline";
import SideMenu from "./components/SideMenu";
import {Route, RouteComponentProps, Switch} from "react-router";
import {BrowserRouter} from "react-router-dom";

import Dashboard from './components/Dashboard';
import Main from './components/Main';
import ScrollToTop from "./components/ScrollTop";

const theme = createMuiTheme({
  palette: {
    secondary: {
      main: blueGrey[500]
    },
    primary: {
      main: blueGrey[800]
    }
  },
  typography: {
    useNextVariants: true,
  }
});

const drawerWidth = 110;
const styles = (theme: any) => ({
  root: {
    display: 'flex',
  },
  drawer: {
    [theme.breakpoints.up('sm')]: {
      width: drawerWidth,
      flexShrink: 0,
    },
  },
  appBar: {
    marginLeft: drawerWidth,
    [theme.breakpoints.up('sm')]: {
      width: `calc(100% - ${drawerWidth}px)`,
    },
  },
  menuButton: {
    marginRight: 20,
    [theme.breakpoints.up('sm')]: {
      display: 'none',
    },
  },
  toolbar: theme.mixins.toolbar,
  drawerPaper: {
    width: drawerWidth,
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing.unit * 3,
  },
});


type AppProps = {
  classes: any,
}

type AppState = {
  sectionName: string
  mobileOpen: boolean,
}

class App extends Component<AppProps, AppState> {
  state = {
    sectionName: "None",
    mobileOpen: false,
    redirectUrl: undefined,
  };

  handleDrawerToggle = () => {
    this.setState(state => ({mobileOpen: !state.mobileOpen}));
  };

  handleSetSectionName = (name: string) => {
    if (this.state.sectionName !== name) {
      this.setState({
        sectionName: name,
      });
    }
  };

  render() {
    const {classes} = this.props;
    return (
      <MuiThemeProvider theme={theme}>
        <div>
          <BrowserRouter>
            <CssBaseline/>
            <AppBar className={classes.appBar} position="sticky">
              <Toolbar>
                <IconButton
                  color="inherit"
                  aria-label="Open drawer"
                  onClick={this.handleDrawerToggle}
                  className={classes.menuButton}>
                  <MenuIcon/>
                </IconButton>
                <Typography variant="h6" color="inherit" noWrap>
                  {this.state.sectionName}
                </Typography>
              </Toolbar>
            </AppBar>
            <SideMenu
              classes={classes}
              theme={theme}
              mobileOpen={this.state.mobileOpen}
              handleDrawerToggle={this.handleDrawerToggle}
            />
            <ScrollToTop>
              <Switch>
                <main className={classes.content}>
                  <Route
                    exact
                    path='/'
                    component={(props: RouteComponentProps) => (
                      <Main {...props} handleSetSectionName={this.handleSetSectionName}/>
                    )}
                  />
                  <Route
                    exact
                    path='/dashboard'
                    component={(props: RouteComponentProps) => (
                      <Dashboard {...props} handleSetSectionName={this.handleSetSectionName}/>
                    )}/>
                </main>
              </Switch>
            </ScrollToTop>
          </BrowserRouter>
        </div>
      </MuiThemeProvider>
    );
  }
}

export default withStyles(styles)(App);
