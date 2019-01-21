<template lang="pug">
  div
    small.text-muted 
      | {{ timestamp | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }} 
      | at {{ timestamp | moment('timezone', 'Asia/Jakarta').format('HH:mm') }}
</template>

<script>
export default {
  name: 'Timestamp',
  props: ['timestamp'],
}
</script>