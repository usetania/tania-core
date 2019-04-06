<template lang="pug">
  footer#footer(role="footer")
    .wrapper.b-t.bg-light
      span.pull-right Tania 1.7.0
        a.m-l-sm.text-muted(href="")
          i.fas.fa-long-arrow-up
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
}
</style>
