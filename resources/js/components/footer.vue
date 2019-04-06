<template lang="pug">
  footer#footer(role="footer")
    .wrapper.bg-light
      | Tania 1.7.0.
      |
      translate Made for the &#9829; of plants &copy; 2019 Copyright.
</template>

<script>
export default {
  name: 'AppFooterComponent',
};
</script>

<style lang="scss" scoped>
#footer {
  position: absolute;
  bottom: 0;
  width: 100%;

  .wrapper {
    height: 60px;
    padding: 10px 20px;
  }

  span {
    font-size: 0.9em;
  }
}

@media (min-width: 768px) {
  #footer {
    .wrapper {
      height: 50px;
    }
  }
}
</style>
