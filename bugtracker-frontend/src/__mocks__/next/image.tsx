import React, { CSSProperties } from "react";

jest.mock("next/image", () => ({
  __esModule: true,
  default: ({
    src,
    alt,
    fill,
    ...props
  }: {
    src: string;
    alt: string;
    fill?: boolean;
  }) => {
    const style: CSSProperties = fill
      ? { position: "absolute", top: 0, left: 0, bottom: 0, right: 0 }
      : {};
    return <img src={src} alt={alt} style={style} {...props} />;
  },
}));
