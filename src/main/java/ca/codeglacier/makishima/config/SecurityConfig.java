package ca.codeglacier.makishima.config;

import ca.codeglacier.makishima.auth.DiscordLoginSuccessHandler;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.csrf.CookieCsrfTokenRepository;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    private final DiscordLoginSuccessHandler successHandler;

    public SecurityConfig(DiscordLoginSuccessHandler successHandler) {
        this.successHandler = successHandler;
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        return http
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/user").authenticated()
                        .anyRequest().permitAll())
                .oauth2Login(oauth -> oauth
                        .successHandler(successHandler))
                .logout(logout -> logout
                        .logoutSuccessUrl("/")
                        .clearAuthentication(true))
                .sessionManagement(session -> session
                        .sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED)
                        .sessionFixation().changeSessionId())
                .build();
    }

}
