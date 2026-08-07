import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";

import Greeting from "../Greeting.vue";

describe("Greeting", () => {
  it("renders the provided name", () => {
    const wrapper = mount(Greeting, { props: { name: "World" } });
    expect(wrapper.text()).toContain("Hello, World!");
  });
});
