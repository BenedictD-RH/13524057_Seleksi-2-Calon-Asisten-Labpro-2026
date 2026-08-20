// Allow global side-effect CSS imports (e.g., import "./App.css")
declare module "*.css";

// Allow CSS Modules if you use them (e.g., import styles from "./App.module.css")
declare module "*.module.css" {
  const classes: { readonly [key: string]: string };
  export default classes;
}
