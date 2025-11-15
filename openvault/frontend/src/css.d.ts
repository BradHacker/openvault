declare module 'csstype' {
  interface Properties {
    [index: `--${string}`]: string | number | undefined;
  }
}
