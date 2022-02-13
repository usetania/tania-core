<template lang="pug">
  div(v-if="description.length > 0")
    p.small(v-if="clicked")
      span {{ description }}
    p
      a(@click="toggle()")
        small.text(v-if="clicked")
          translate Close details
          i.fa.fa-angle-up
        small.text(v-else)
          translate Read details
          i.fa.fa-angle-down
</template>

<script>
export default {
  name: 'MoreDetail',
  props: ['description'],
  data() {
    return {
      clicked: false,
    };
  },
  methods: {
    toggle() {
      this.clicked = !this.clicked;
    },
  },
};
</script>
