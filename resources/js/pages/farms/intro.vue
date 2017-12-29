<template lang="pug">
  .container.init.col-md-4.col-md-offset-4
    a.navbar-brand.block.m-b.m-t.text-center
      img(src="../../../images/logobig.png")
    FarmCreate(:intro="true")
</template>

<script>
export default {
  name: 'FarmIntro',
  components: {
    FarmCreate: () => import('./create')
  }
}
</script>

