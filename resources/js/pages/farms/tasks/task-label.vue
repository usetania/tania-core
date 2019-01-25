<template lang="pug">
  div
    span
      span(v-if="type == 'CATEGORY'")
        span.label.label-pestcontrol(v-if="task.category == 'PESTCONTROL'") PEST CONTROL
        span.label.label-sanitation(v-if="task.category == 'SANITATION'") SANITATION
        span.label.label-area(v-if="task.category == 'AREA'") AREA
        span.label.label-reservoir(v-if="task.category == 'RESERVOIR'") RESERVOIR
        span.label.label-safety(v-if="task.category == 'SAFETY'") SAFETY
        span.label.label-finance(v-if="task.category == 'FINANCE'") FINANCE
        span.label.label-crop(v-if="task.category == 'CROP'") CROP
        span.label.label-water(v-if="task.category == 'WATER'") WATER
        span.label.label-inventory(v-if="task.category == 'INVENTORY'") INVENTORY
        span.label.label-general(v-if="task.category == 'GENERAL'") GENERAL
        span.label.label-nutrient(v-if="task.category == 'NUTRIENT'") NUTRIENT
      span(v-if="type == 'PRIORITY'")
        span.status.status-urgent(v-if="task.priority == 'URGENT'") URGENT
        span.status.status-normal(v-if="task.priority == 'NORMAL'") NORMAL
</template>

<script>
export default {
  name: 'TaskLabel',
  props: ['type', 'task'],
}
</script>