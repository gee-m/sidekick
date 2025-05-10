# Life Tracking System

The goal is to make tracking "stupidly easy" - as easy as pressing a button or sending a quick message, or maybe even snapping a pic. Some questions are things that should be asked daily (e.g., how you feel), some should be inputted by the user but may need a reminder if not inputted during the day.

## Project Overview

A comprehensive life tracking system designed to collect, correlate, and analyze various aspects of daily life to identify patterns affecting happiness and wellbeing. Inspired by Bryan Johnson's tracking work and Andrej Karpathy's sleep tracking analysis.

### Core Philosophy

> "My sleep scores correlate strongly with the quality of work I am able to do that day. When my score is low, I lack agency, I lack courage, I lack creativity, I'm simply tired. When my sleep score is high, I can power through anything." - Andrej Karpathy

This system aims to identify similar correlations across multiple life domains to optimize wellbeing through data-driven insights.

## Core Objectives

1. **Low Friction Data Collection** - Make tracking "stupidly easy" through mobile-optimized interfaces
2. **Comprehensive Life Domains** - Track sleep, nutrition, exercise, work, relationships, and more
3. **Meaningful Correlations** - Identify connections between behaviors and wellbeing outcomes
4. **Actionable Insights** - Generate practical recommendations based on personal data patterns

## Domain Model

### DailyRecord
- `date`: Date - The date of this record
- `subjective`: SubjectiveMetrics - Subjective wellbeing indicators
- `sleep`: SleepData - Sleep metrics including Oura data
- `nutrition`: NutritionData - Food, drink, and meal timing data
- `activities`: ActivityData - Exercise, nature time, and other activities
- `relationships`: RelationshipData - Social and relationship interactions
- `supplements`: SupplementData - Supplement and medication tracking
- `timeUsage`: TimeUsageData - How time was allocated across activities
- `selfCare`: SelfCareData - Self-care activities and tracking
- `work`: WorkData - Work-related metrics and productivity
- `derivedMetrics`: DerivedMetrics - Automatically calculated metrics (streaks, days since, etc.)

### QuickAction
- `id`: String - Unique identifier
- `type`: Enum(SUPPLEMENT, MEAL, ACTIVITY, SELFCARE, etc.) - Type of action
- `label`: String - Display text for the button
- `payload`: JSON - Data to be submitted when tapped
- `icon`: String - Icon identifier
- `usageCount`: Number - How many times action has been used
- `lastUsed`: DateTime - When action was last used
- `contextTags`: [String] - When this action is most relevant (morning, evening, etc.)

### DerivedMetrics
- `streaks`: Map<ActivityType, Number> - Consecutive days of activities/habits
- `daysSinceEvents`: Map<EventType, Number> - Days since various tracked events
  - Contains all "days since" metrics (work, free day, holiday, fast food, gluten, 
    family contact, exercise, supplements, nature time, animal interaction, etc.)
- `consistencyScores`: Map<MetricType, Score> - Consistency of various metrics
- `correlationInsights`: [CorrelationInsight] - Automatically detected correlations

### SubjectiveMetrics
- `morningRestedness`: Score(1-10) - How rested upon waking
- `morningHappiness`: Score(1-10) - Happiness level in morning
- `morningAnxiety`: Score(1-10) - Anxiety level in morning
- `eveningHappiness`: Score(1-10) - Happiness level in evening
- `eveningAnxiety`: Score(1-10) - Anxiety level in evening
- `productivityFeeling`: Score(1-10) - Subjective productivity rating
- `lifeAdminPressure`: Score(1-10) - Pressure from life admin tasks

### SleepData
- `ouraData`: OuraMetrics - Data imported from Oura API
  - `sleepScore`: Number - Overall sleep score
  - `deepSleep`: Duration - Time in deep sleep
  - `remSleep`: Duration - Time in REM sleep
  - `lightSleep`: Duration - Time in light sleep
  - `awake`: Duration - Time awake during sleep period
  - `totalSleep`: Duration - Total sleep duration
  - `restingHeartRate`: Number - RHR during sleep
  - `heartRateVariability`: Number - HRV during sleep
  - `respiratoryRate`: Number - Breathing rate during sleep
  - `bodyTemperature`: Number - Body temperature deviation
- `bedTime`: Time - Time went to bed
- `wakeTime`: Time - Time woke up
- `bedTimeConsistency7d`: Duration - 7-day consistency of bedtime
- `bedTimeConsistency30d`: Duration - 30-day consistency of bedtime
- `wakeTimeConsistency7d`: Duration - 7-day consistency of wake time
- `wakeTimeConsistency30d`: Duration - 30-day consistency of wake time
- `hoursBeforeSleepLastMeal`: Number - Hours between last meal and sleep

### NutritionData
- `meals`: [Meal] - List of meals consumed
  - `mealType`: Enum(BREAKFAST, LUNCH, DINNER, SNACK)
  - `time`: Time - When meal was consumed
  - `category`: [FoodCategory] - Types of food consumed
  - `satiety`: Score(1-10) - How filling the meal was (10 = overate)
  - `glutenConsumed`: Boolean - Whether meal contained gluten
  - `ingredients`: [String] - Key ingredients in the meal
- `fastFoodConsumed`: Boolean - Whether fast food was consumed
- `alcoholConsumed`: Boolean - Whether alcohol was consumed
- `caffeineConsumed`: Boolean - Whether caffeine was consumed
- `waterIntake`: Number - Estimated water consumption in ml

### ActivityData
- `stepCount`: Number - Number of steps taken
- `exercisePerformed`: Boolean - Whether exercise was performed
- `exerciseType`: [ExerciseType] - Types of exercise performed
- `exerciseDuration`: Duration - Duration of exercise
- `outdoorTime`: Duration - Time spent outdoors
- `natureTime`: Duration - Time spent in nature

### RelationshipData
- `partnerProximity`: Boolean - Whether partner was physically present
- `partnerInteraction`: Score(1-10) - Quality of partner interactions
- `intimacyWithPartner`: Boolean - Whether intimacy occurred with partner
- `conflictWithPartner`: Boolean - Whether conflict occurred with partner
- `conflictIntensity`: Score(1-10) - Intensity of any conflict
- `familyContact`: Boolean - Whether family contact occurred
- `socialInteraction`: Boolean - Whether social interaction occurred
- `animalInteraction`: Boolean - Whether interaction with pets/animals occurred

### SupplementData
- `supplementsTaken`: [Supplement] - Supplements consumed
  - `name`: String - Name of supplement
  - `dose`: String - Dose of supplement
  - `time`: Time - Time supplement was taken
- `creatineAmount`: Number - Amount of creatine consumed (mg)
- `dutasterideTaken`: Boolean - Whether dutasteride was taken
- `minoxidilApplied`: Boolean - Whether minoxidil was applied
- `otherMedications`: [Medication] - Other medications taken

### TimeUsageData
- `workScreenTime`: Duration - Screen time for work
- `personalComputerTime`: Duration - Personal computer usage time
- `phoneScreenTime`: Duration - Phone screen time
- `phoneAppBreakdown`: Map<App, Duration> - Breakdown of app usage
  - Including Instagram, YouTube, etc.
- `workPerformed`: Boolean - Whether work was performed
- `workProductivity`: Score(1-10) - Subjective work productivity

### SelfCareData
- `selfCareActivities`: [SelfCareActivity] - Self-care activities performed
  - `type`: Enum(HAIRCUT, BATH, MASSAGE, SKINCARE, MEDITATION, OTHER)
  - `duration`: Duration - Duration of activity
  - `satisfaction`: Score(1-10) - Satisfaction with activity

### WorkData
- `workDay`: Boolean - Whether it was a work day
- `workHours`: Duration - Hours worked
- `meetingTime`: Duration - Time spent in meetings
- `focusTime`: Duration - Time spent in focused work
- `workSatisfaction`: Score(1-10) - Satisfaction with work accomplished

## Implementation Phases

### Phase -1: LLM-Based Conversational Tracking
- **Daily Check-in Structure**:
  - **Morning Check-in**:
    - "How rested do you feel today? (1-10)"
    - "How happy do you feel this morning? (1-10)"
    - "How anxious do you feel this morning? (1-10)"
    - "What time did you go to bed and wake up?"
    - "Any supplements taken this morning?"
  
  - **Evening Check-in**:
    - "How happy do you feel tonight? (1-10)"
    - "How anxious do you feel tonight? (1-10)"
    - "How productive were you today? (1-10)"
    - "Did you work today or was it a free day?"
    - "What did you eat today? Any fast food or gluten?"
    - "Did you exercise today? Nature time? Step count?"
    - "Any supplements taken today? Creatine?"
    - "How much time did you spend on screens today?"
    - "How was your relationship with your partner today?"
    - "Any self-care activities today?"
    - "How much pressure do you feel from life admin? (1-10)"

- **LLM Prompt Template**:
  ```
  You are a life tracking assistant. Your job is to collect daily data points about my wellbeing,
  activities, and habits, then format this data in a structured JSON format for later analysis.
  
  Ask me the appropriate morning or evening questions based on the time of day.
  After I respond to all questions, output a JSON object with all collected data.
  
  This JSON object should follow this structure:
  {
    "date": "YYYY-MM-DD",
    "timeOfDay": "morning|evening",
    "subjective": {
      "restedness": 1-10,
      "happiness": 1-10,
      "anxiety": 1-10,
      "productivity": 1-10,
      "lifeAdminPressure": 1-10
    },
    "sleep": {
      "bedTime": "HH:MM",
      "wakeTime": "HH:MM"
    },
    ...other relevant fields based on answers...
  }
  
  Store this data incrementally, so we can build up a complete daily record.
  Remind me if I miss a check-in.
  ```

- Manual data collection through structured LLM conversations
- JSON data storage for easy migration to later phases
- Refinement of metrics and tracking categories
- Identification of key correlations and patterns
- Development of tracking habits before system implementation

### Phase 1: Mobile-First Basic System
- Go backend with HTMX frontend optimized for mobile
- Progressive Web App (PWA) for home screen installation 
- Quick-tap interfaces for common inputs
- **Dynamic quick-action system** for frequently used inputs
- PostgreSQL data storage
- Basic visualization of recent trends
- Full Oura Ring API integration
- **Robust offline capabilities**:
  - Complete offline functionality for all data entry
  - Background synchronization when connection restored
  - Conflict resolution for offline changes
  - Local storage with encryption

### Phase 2: Integration & Automation
- Expanded API integrations (where available)
- Automated tracking for objective metrics
- Smart defaults and pattern recognition
- Notification system with actionable inputs
- Reminder system for daily inputs

### Phase 3: Analysis & Advanced Features
- Statistical analysis for identifying correlations
- Comprehensive dashboard across devices
- Advanced pattern recognition
- Recommendation engine
- Data export/sharing capabilities

## Statistical Methodology

### Data Analysis Approach

#### Handling Missing Data
- **Interpolation Methods**:
  - Linear interpolation for continuous metrics (like sleep duration)
  - Last observation carried forward for categorical data
  - Multiple imputation for critical variables when appropriate
- **Missing Data Visualization**: Clear indication of interpolated/missing data in visualizations
- **Tracking Consistency Metrics**: Track and visualize data collection completeness

#### Correlation Analysis
- **Lag Analysis**: Examine effects across different time frames (same-day, next-day, multi-day)
- **Rolling Windows**: Use 3-day, 7-day, and 30-day windows to identify longer-term patterns
- **Multivariate Analysis**: Control for confounding variables when establishing correlations
  - Example: Sleep quality vs. mood, controlling for exercise and nutrition
- **Causality Testing**: Granger causality tests to suggest potential causal relationships

#### Confounding Variable Management
- **Factor Analysis**: Identify clusters of related variables
- **Stratification**: Analyze correlations within subsets of similar days
- **Control Variables**: Maintain awareness of external factors (weather, seasons, work demands)
- **Natural Experiments**: Leverage unplanned disruptions (travel, illness) to isolate variables

#### Insight Generation
- **Hypothesis Testing**: Formal testing of correlations with statistical significance
- **Anomaly Detection**: Identify outlier days and examine contributing factors
- **Pattern Recognition**: Machine learning techniques to identify complex patterns
  - Decision trees for identifying key decision points
  - Clustering for identifying similar "types of days"
- **Personalized Baselines**: Establish individual baselines rather than population norms

### Visualization Strategies
- **Correlation Heatmaps**: Visualize strength of relationships between variables
- **Time Series Analysis**: Track trends, seasonality, and cyclical patterns
- **Comparative Analysis**: Side-by-side comparison of different time periods
- **Causal Diagrams**: Visualize potential causal relationships between variables

## Technical Architecture

### Backend
- **Language**: Go
- **Web Framework**: [Appropriate Go web framework]
- **Database**: PostgreSQL
- **API**: RESTful endpoints + GraphQL (optional)
- **Authentication**: OAuth2 / JWT
- **Deployment**: render.com infrastructure

### Frontend
- **Primary Technology**: HTMX + Alpine.js
- **CSS Framework**: TailwindCSS
- **Mobile Optimization**: PWA with offline capabilities
  - Service workers for offline data collection
  - IndexedDB for local data storage
  - Background sync API for data transmission when online
  - Conflict resolution for simultaneous online/offline changes
- **Data Visualization**: Chart.js or D3.js
- **Quick Actions Engine**: Client/server system for dynamic action generation based on usage patterns

### API Integrations
- **Oura Ring API**: https://cloud.ouraring.com/docs/
  - Authentication via OAuth2
  - Daily syncing of sleep, activity, and readiness data
  - Historical data import
- **Future Integrations**:
  - Phone screen time APIs
  - Calendar integration for work/vacation detection
  - Fitness tracker APIs

## Oura API Integration Specifications

### Authentication
- OAuth2 flow for authorization
- Token storage with refresh capabilities
- User permission for accessing:
  - Sleep data
  - Activity data
  - Readiness data

### Data Retrieval
- Daily synchronization job for latest data
- Endpoints to use:
  - `/v2/usercollection/sleep`
  - `/v2/usercollection/daily_activity`
  - `/v2/usercollection/daily_readiness`
- Key metrics to extract:
  - Sleep scores and stages
  - HRV and resting heart rate
  - Body temperature
  - Activity levels
  - Readiness score

### Processing
- Calculate sleep consistency metrics
- Normalize data for consistent tracking
- Generate derived metrics (e.g., consistency scores)

## User Experience Specifications

### Mobile Input Principles
- **One-Thumb Operation**: Design for single-handed use
- **Tap Over Type**: Prefer tap interactions over keyboard input
- **Smart Defaults**: Pre-fill with likely values based on patterns
- **Progressive Disclosure**: Hide complexity, reveal as needed
- **Batch Inputs**: Group related inputs to reduce context switching
- **Micro-Interactions**: Quick inputs throughout day vs. one long form
- **Dynamic Quick Actions**: Automatically generate buttons for recently/frequently used inputs

### Key Interfaces
- **Morning Check-In**:
  - Quick subjective scores (restedness, happiness, anxiety)
  - Confirm/adjust sleep timings from Oura
  - Previous day recap for missing data
  - Dynamic quick-action buttons for common morning activities/supplements

- **Evening Check-In**:
  - Meal tracking interface with quick categories
  - Supplement tracking with favorites system
  - Relationship and interaction tracking
  - Self-care activity logging
  - Work productivity assessment
  - Dynamic quick-action buttons for common evening activities

- **Ad-Hoc Tracking**:
  - Quick add for activities, meals, supplements
  - Photo input option for meals
  - Voice note capability for qualitative data
  - Recent/frequent items as tappable buttons (top 10-20)
  - Contextual suggestions based on time of day and patterns

- **Action Logging System**:
  - Reverse-chronological timeline of all logged actions
  - Daily summary of everything logged
  - Missing data alerts for categories not yet logged
  - Easily accessible with simple swipe gesture from main interface
  - Option to edit or delete recently logged items

## Next Steps

1. **Immediate Data Collection Setup**
   - Create Typeform templates for morning and evening check-ins
   - Set up Google Sheets integration for data storage
   - Develop simple WhatsApp bot for ad-hoc logging
   - Begin collecting data through these channels

2. **Oura Integration**
   - Set up Oura API authentication
   - Develop basic data import functionality
   - Create data normalization processes

3. **PWA Development**
   - Set up repository structure
   - Implement core domain models in Go
   - Create database schema on render.com
   - Develop mobile-first UI with HTMX
   - Build quick-actions engine

4. **Analysis Capabilities**
   - Implement basic correlation analysis
   - Develop visualization dashboard
   - Create daily/weekly report generation

## File Structure

```
/cmd
  /api        # API server entry point
  /migrations # Database migrations
/internal
  /domain     # Domain model definitions
  /service    # Business logic implementation
  /repository # Data storage interfaces
  /api        # API handlers
  /integration # Third-party integrations
    /oura     # Oura Ring API integration
  /analysis   # Data analysis capabilities
/web
  /templates  # HTMX templates
  /static     # Static assets
  /js         # Client-side JavaScript
/scripts      # Utility scripts
/docs         # Documentation
```

## Limitations and Considerations

### Data Collection Challenges
- **Tracking Fatigue Risk**: The comprehensive nature of metrics may lead to inconsistent tracking
- **Solution**: Implement a tiered approach with "core" vs "optional" metrics
- **Consistency Challenges**: Multiple input methods could create data fragmentation
- **Solution**: Define strict data schemas and normalization processes

### Privacy and Security
- **Sensitive Data**: System contains personal health and relationship information
- **Data Protection**: Implement appropriate encryption and access controls
- **Data Ownership**: All data remains user-owned with complete export capabilities
- **Backup Strategy**: Regular automated backups with versioning

### Statistical Analysis Considerations
- **Confounding Variables**: Many life factors correlate without causation
- **Missing Data Handling**: Statistical methods to account for tracking gaps
- **Insight Validation**: Mechanisms to validate and test correlation hypotheses

### Technical Considerations
- **Offline Capability**: PWA must function without internet connection
- **Data Synchronization**: Conflict resolution for offline data
- **Progressive Enhancement**: Core functionality works without advanced features
- **Testing Strategy**: Automated tests for data integrity and UI functionality

### User Experience Challenges
- **Consistency vs. Convenience**: Balance between detailed data and tracking ease
- **Notification Management**: Avoid alert fatigue while maintaining tracking consistency
- **Learning Curve**: Interface complexity increases with feature richness
- **Solution**: Implement progressive disclosure of advanced features

## Original Tracking Requirements

The following is the original set of tracking requirements that informed this specification:

```
Together I want us to build something that will let me tell it how I feel, what I take, how I've slept, with a general formula of trying to figure out: what i took yesterday, how i felt, how much i exercised, how much i worked, how much stress = how much sleep + how good i feel the next day

I want to track things like, in order to import data somewhere useful so we can analyze it:
* All Oura Ring data on my sleep, activities, readiness, step count
* How do I feel on wake up (rested score, happiness score, anxiety score)
* How did I feel today (Rested score, happiness score, anxiety score)
* How productive I felt at work today
* Whether I worked, or am off (weekends count as off)
* Auto track, how long it's been:
   * since I worked
   * since free day
   * since at least 1w holiday
   * since fast food
   * since gluten, maybe other potential toxins
   * since called family
   * since exercise, since at least 4k steps
   * since XX supplement
   * since caffeine
   * since started taking dutasteride, minoxidil
   * since I was last in nature
   * since interacted with animals (cat, dog, etc)
   * since I worked out (eg push ups count)
   * since partner was next to me
   * since slept with partner
   * since (big) fight with partner
   * (some of this stuff means I also need to track daily, like whether i was in nature or not! make sure to understand that as well)
* How consistent was the time in going to bed today, and 7d, 30d average
* How consistent was wake up time today, and 7d, 30d average
* How many hours did I eat before sleep?
* What type of food/drinks i ate (discrete categories, eg fast food, salad, mexican, sushi, etc...)
   * Whether there was gluten (or other potential toxins) (discrete categories)
   * Maybe instead listing ingredients also? can be auto determined eg if burger
   * for breakfast (if had any), lunch, dinner
   * how much satiety (1-10), 10 = overate
* How much creatine I took, when
* What supplements I took, when
* How much screen time (work comp, personal comp, and android: instagram, phone, youtube)
* How much time spent outside
* How much pressure do I feel from life admin
* How today's relationship with partner
```